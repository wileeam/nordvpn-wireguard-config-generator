package gen

import (
	"math"
	"sort"
	"testing"

	"github.com/mustafachyi/nordvpn-wireguard-config-generator/internal/structs"
)

// ---------------------------------------------------------------------------
// checkVer
// ---------------------------------------------------------------------------

func TestCheckVer(t *testing.T) {
	cases := []struct {
		v    string
		want bool
	}{
		{"2.1.0", true},
		{"2.1", true},
		{"2.2.0", true},
		{"2.0.0", false},
		{"2.0", false},
		{"1.9.9", false},
		{"3.0.0", true},
		{"10.0.0", true},
		{"0.0.0", false},
		{"", false},
		{"2", false},
		{"a.b.c", false},
		{"2.1.abc", true},  // minor parsed up to first non-digit
		{"2.10.0", true},
		{"2.0.9", false},
	}

	for _, tc := range cases {
		got := checkVer(tc.v)
		if got != tc.want {
			t.Errorf("checkVer(%q) = %v, want %v", tc.v, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// haversine
// ---------------------------------------------------------------------------

func TestHaversine(t *testing.T) {
	t.Run("same_point_is_zero", func(t *testing.T) {
		d := haversine(0, 0, 0, 0)
		if d != 0 {
			t.Errorf("expected 0, got %f", d)
		}
	})

	t.Run("ny_to_london", func(t *testing.T) {
		d := haversine(40.7128, -74.0060, 51.5074, -0.1278)
		const expected = 5570.0
		const tol = 0.01
		if math.Abs(d-expected)/expected > tol {
			t.Errorf("NY→London: expected ~%f km, got %f km", expected, d)
		}
	})

	t.Run("symmetric", func(t *testing.T) {
		d1 := haversine(10, 20, 30, 40)
		d2 := haversine(30, 40, 10, 20)
		if math.Abs(d1-d2) > 1e-9 {
			t.Errorf("haversine not symmetric: %f vs %f", d1, d2)
		}
	})

	t.Run("antipodal", func(t *testing.T) {
		d := haversine(0, 0, 0, 180)
		expected := math.Pi * 6371
		if math.Abs(d-expected)/expected > 0.001 {
			t.Errorf("antipodal: expected ~%f, got %f", expected, d)
		}
	})

	t.Run("positive_result", func(t *testing.T) {
		d := haversine(0, 0, 10, 10)
		if d <= 0 {
			t.Errorf("expected positive distance, got %f", d)
		}
	})
}

// ---------------------------------------------------------------------------
// byLoadDist sort
// ---------------------------------------------------------------------------

func TestByLoadDist(t *testing.T) {
	servers := byLoadDist{
		{Load: 50, Dist: 100},
		{Load: 10, Dist: 500},
		{Load: 10, Dist: 200},
		{Load: 80, Dist: 50},
	}
	sort.Sort(servers)

	// After sort: load 10 dist 200, load 10 dist 500, load 50 dist 100, load 80 dist 50
	if servers[0].Load != 10 || servers[0].Dist != 200 {
		t.Errorf("expected first: load=10 dist=200, got load=%d dist=%f", servers[0].Load, servers[0].Dist)
	}
	if servers[1].Load != 10 || servers[1].Dist != 500 {
		t.Errorf("expected second: load=10 dist=500, got load=%d dist=%f", servers[1].Load, servers[1].Dist)
	}
	if servers[2].Load != 50 {
		t.Errorf("expected third load=50, got %d", servers[2].Load)
	}
	if servers[3].Load != 80 {
		t.Errorf("expected fourth load=80, got %d", servers[3].Load)
	}
}

func TestByLoadDistSameLoadSameDist(t *testing.T) {
	servers := byLoadDist{
		{Load: 30, Dist: 100},
		{Load: 30, Dist: 100},
	}
	sort.Sort(servers)
	// No panic, order is stable (or not - just ensure no crash)
}

// ---------------------------------------------------------------------------
// parse
// ---------------------------------------------------------------------------

func makeRawServer() structs.ApiServer {
	return structs.ApiServer{
		Name:     "US#1",
		Hostname: "us1.nordvpn.com",
		Station:  "1.2.3.4",
		Load:     25,
		Locations: []structs.ApiLocation{
			{
				Lat: 40.7,
				Lon: -74.0,
				Country: struct {
					Name string `json:"name"`
					Code string `json:"code"`
					City struct {
						Name string `json:"name"`
					} `json:"city"`
				}{
					Name: "United States",
					Code: "US",
					City: struct {
						Name string `json:"name"`
					}{Name: "New York"},
				},
			},
		},
		Tech: []structs.ApiTech{
			{
				ID: "wireguard_udp",
				Meta: []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				}{
					{Name: "public_key", Value: "testpublickey="},
				},
			},
		},
		Specs: []structs.ApiSpec{
			{
				ID: "version",
				Values: []struct {
					Value string `json:"value"`
				}{
					{Value: "2.1.0"},
				},
			},
		},
	}
}

func TestParseValid(t *testing.T) {
	raw := makeRawServer()
	s := parse(&raw, 0, 0)
	if s == nil {
		t.Fatal("expected non-nil Server")
	}
	if s.Name != "US#1" {
		t.Errorf("Name: got %q, want %q", s.Name, "US#1")
	}
	if s.Host != "us1.nordvpn.com" {
		t.Errorf("Host: got %q", s.Host)
	}
	if s.IP != "1.2.3.4" {
		t.Errorf("IP: got %q", s.IP)
	}
	if s.Load != 25 {
		t.Errorf("Load: got %d", s.Load)
	}
	if s.Country != "United States" {
		t.Errorf("Country: got %q", s.Country)
	}
	if s.Code != "us" {
		t.Errorf("Code: got %q, want lowercase 'us'", s.Code)
	}
	if s.City != "New York" {
		t.Errorf("City: got %q", s.City)
	}
	if s.PubK != "testpublickey=" {
		t.Errorf("PubK: got %q", s.PubK)
	}
}

func TestParseNoLocations(t *testing.T) {
	raw := makeRawServer()
	raw.Locations = nil
	if parse(&raw, 0, 0) != nil {
		t.Error("expected nil for missing locations")
	}
}

func TestParseEmptyLocations(t *testing.T) {
	raw := makeRawServer()
	raw.Locations = []structs.ApiLocation{}
	if parse(&raw, 0, 0) != nil {
		t.Error("expected nil for empty locations")
	}
}

func TestParseOldVersion(t *testing.T) {
	raw := makeRawServer()
	raw.Specs[0].Values[0].Value = "2.0.0"
	if parse(&raw, 0, 0) != nil {
		t.Error("expected nil for old version")
	}
}

func TestParseNoWireGuardTech(t *testing.T) {
	raw := makeRawServer()
	raw.Tech = []structs.ApiTech{{ID: "openvpn_tcp"}}
	if parse(&raw, 0, 0) != nil {
		t.Error("expected nil when no wireguard_udp tech")
	}
}

func TestParseNoPublicKey(t *testing.T) {
	raw := makeRawServer()
	raw.Tech = []structs.ApiTech{
		{ID: "wireguard_udp", Meta: nil},
	}
	if parse(&raw, 0, 0) != nil {
		t.Error("expected nil when public_key is missing")
	}
}

func TestParseNoSpecs(t *testing.T) {
	raw := makeRawServer()
	raw.Specs = nil
	// default version "0.0.0" fails version check
	if parse(&raw, 0, 0) != nil {
		t.Error("expected nil when no version spec")
	}
}

func TestParseCountryCodeLowercased(t *testing.T) {
	raw := makeRawServer()
	raw.Locations[0].Country.Code = "DE"
	s := parse(&raw, 0, 0)
	if s == nil {
		t.Fatal("expected non-nil")
	}
	if s.Code != "de" {
		t.Errorf("Code should be lowercase: got %q", s.Code)
	}
}

// ---------------------------------------------------------------------------
// process (deduplication)
// ---------------------------------------------------------------------------

func TestProcessDeduplication(t *testing.T) {
	raw1 := makeRawServer()
	raw2 := makeRawServer() // same name "US#1"
	raw3 := makeRawServer()
	raw3.Name = "US#2"

	result, rejected := process([]structs.ApiServer{raw1, raw2, raw3}, 0, 0)

	if len(result) != 2 {
		t.Errorf("expected 2 unique servers, got %d", len(result))
	}
	// rejected = 3 raw - 2 unique = 1
	if rejected != 1 {
		t.Errorf("expected 1 rejected, got %d", rejected)
	}
}

func TestProcessEmpty(t *testing.T) {
	result, rejected := process(nil, 0, 0)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil input")
	}
	if rejected != 0 {
		t.Errorf("expected 0 rejected for nil input")
	}
}

func TestProcessFiltersInvalid(t *testing.T) {
	bad := makeRawServer()
	bad.Locations = nil

	good := makeRawServer()
	good.Name = "DE#1"

	result, rejected := process([]structs.ApiServer{bad, good}, 0, 0)
	if len(result) != 1 {
		t.Errorf("expected 1 valid server, got %d", len(result))
	}
	if result[0].Name != "DE#1" {
		t.Errorf("expected DE#1, got %s", result[0].Name)
	}
	_ = rejected
}
