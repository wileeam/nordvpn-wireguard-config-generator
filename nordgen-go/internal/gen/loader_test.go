package gen

import (
	"testing"

	"github.com/mustafachyi/nordvpn-wireguard-config-generator/internal/structs"
)

// TestCheckVer validates the version checking function
func TestCheckVer(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"valid 2.1", "2.1.0", true},
		{"valid 2.2", "2.2.0", true},
		{"valid 3.0", "3.0.0", true},
		{"valid 3.1", "3.1.5", true},
		{"valid 10.5", "10.5.2", true},
		{"invalid 2.0", "2.0.0", false},
		{"invalid 1.9", "1.9.9", false},
		{"invalid 0.5", "0.5.0", false},
		{"invalid 1.0", "1.0", false},
		{"too short", "2", false},
		{"too short no dot", "21", false},
		{"empty", "", false},
		{"no dots", "210", false},
		{"non-numeric major", "a.1.0", false},
		{"non-numeric minor", "2.a.0", false},
		{"edge case 2.1", "2.1", true},
		{"edge case 2.0", "2.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkVer(tt.version)
			if got != tt.want {
				t.Errorf("checkVer(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

// TestHaversine validates the haversine distance calculation
func TestHaversine(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		wantMin  float64
		wantMax  float64
	}{
		{
			name:    "same location",
			lat1:    40.7128,
			lon1:    -74.0060,
			lat2:    40.7128,
			lon2:    -74.0060,
			wantMin: 0,
			wantMax: 0.01,
		},
		{
			name:    "New York to Los Angeles",
			lat1:    40.7128,
			lon1:    -74.0060,
			lat2:    34.0522,
			lon2:    -118.2437,
			wantMin: 3900,
			wantMax: 4000,
		},
		{
			name:    "1 degree longitude at equator",
			lat1:    0,
			lon1:    0,
			lat2:    0,
			lon2:    1,
			wantMin: 110,
			wantMax: 112,
		},
		{
			name:    "1 degree latitude",
			lat1:    0,
			lon1:    0,
			lat2:    1,
			lon2:    0,
			wantMin: 110,
			wantMax: 112,
		},
		{
			name:    "London to Paris",
			lat1:    51.5074,
			lon1:    -0.1278,
			lat2:    48.8566,
			lon2:    2.3522,
			wantMin: 340,
			wantMax: 350,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := haversine(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("haversine(%f, %f, %f, %f) = %f, want between %f and %f",
					tt.lat1, tt.lon1, tt.lat2, tt.lon2, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// TestParse validates the server parsing function
func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		server *structs.ApiServer
		lat    float64
		lon    float64
		want   bool // true if server should be parsed successfully
	}{
		{
			name: "valid server",
			server: &structs.ApiServer{
				Name:     "us1234",
				Hostname: "us1234.nordvpn.com",
				Station:  "192.168.1.1",
				Load:     50,
				Locations: []structs.Location{
					{
						Country: structs.Country{
							Name: "United States",
							Code: "US",
							City: structs.City{Name: "New York"},
						},
						Lat: 40.7128,
						Lon: -74.0060,
					},
				},
				Specs: []structs.Spec{
					{
						ID: "version",
						Values: []structs.SpecVal{
							{Value: "2.1.0"},
						},
					},
				},
				Tech: []structs.Tech{
					{
						ID: "wireguard_udp",
						Meta: []structs.Meta{
							{Name: "public_key", Value: "testkey123="},
						},
					},
				},
			},
			lat:  40.0,
			lon:  -73.0,
			want: true,
		},
		{
			name: "no locations",
			server: &structs.ApiServer{
				Name:      "us1234",
				Locations: []structs.Location{},
				Specs: []structs.Spec{
					{ID: "version", Values: []structs.SpecVal{{Value: "2.1.0"}}},
				},
				Tech: []structs.Tech{
					{ID: "wireguard_udp", Meta: []structs.Meta{{Name: "public_key", Value: "key"}}},
				},
			},
			lat:  40.0,
			lon:  -73.0,
			want: false,
		},
		{
			name: "old version",
			server: &structs.ApiServer{
				Name:     "us1234",
				Hostname: "us1234.nordvpn.com",
				Station:  "192.168.1.1",
				Locations: []structs.Location{
					{
						Country: structs.Country{Name: "US", Code: "US", City: structs.City{Name: "NY"}},
						Lat:     40.0,
						Lon:     -73.0,
					},
				},
				Specs: []structs.Spec{
					{ID: "version", Values: []structs.SpecVal{{Value: "1.0.0"}}},
				},
				Tech: []structs.Tech{
					{ID: "wireguard_udp", Meta: []structs.Meta{{Name: "public_key", Value: "key"}}},
				},
			},
			lat:  40.0,
			lon:  -73.0,
			want: false,
		},
		{
			name: "no public key",
			server: &structs.ApiServer{
				Name:     "us1234",
				Hostname: "us1234.nordvpn.com",
				Station:  "192.168.1.1",
				Locations: []structs.Location{
					{
						Country: structs.Country{Name: "US", Code: "US", City: structs.City{Name: "NY"}},
						Lat:     40.0,
						Lon:     -73.0,
					},
				},
				Specs: []structs.Spec{
					{ID: "version", Values: []structs.SpecVal{{Value: "2.1.0"}}},
				},
				Tech: []structs.Tech{},
			},
			lat:  40.0,
			lon:  -73.0,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(tt.server, tt.lat, tt.lon)
			if (got != nil) != tt.want {
				t.Errorf("parse() returned %v, want success=%v", got, tt.want)
			}

			if got != nil {
				// Validate parsed server fields
				if got.Name != tt.server.Name {
					t.Errorf("parsed server Name = %q, want %q", got.Name, tt.server.Name)
				}
				if got.Host != tt.server.Hostname {
					t.Errorf("parsed server Host = %q, want %q", got.Host, tt.server.Hostname)
				}
				if got.IP != tt.server.Station {
					t.Errorf("parsed server IP = %q, want %q", got.IP, tt.server.Station)
				}
				if got.Dist < 0 {
					t.Error("parsed server Dist should be non-negative")
				}
			}
		})
	}
}

// TestByLoadDist validates the sorting implementation
func TestByLoadDist(t *testing.T) {
	servers := []structs.Server{
		{Name: "high_load_far", Load: 80, Dist: 1000},
		{Name: "low_load_far", Load: 20, Dist: 1000},
		{Name: "low_load_near", Load: 20, Dist: 100},
		{Name: "medium_load_medium", Load: 50, Dist: 500},
	}

	sortable := byLoadDist(servers)

	// Test Len
	if sortable.Len() != 4 {
		t.Errorf("Len() = %d, want 4", sortable.Len())
	}

	// Test Less - should prioritize by load, then distance
	if !sortable.Less(2, 1) { // low_load_near < low_load_far (same load, shorter distance)
		t.Error("Less() should prioritize shorter distance when loads are equal")
	}

	if !sortable.Less(1, 3) { // low_load_far < medium_load_medium (lower load)
		t.Error("Less() should prioritize lower load")
	}

	if !sortable.Less(3, 0) { // medium_load_medium < high_load_far (lower load)
		t.Error("Less() should prioritize lower load")
	}

	// Test Swap
	sortable.Swap(0, 1)
	if servers[0].Name != "low_load_far" || servers[1].Name != "high_load_far" {
		t.Error("Swap() did not swap elements correctly")
	}
}

// TestProcess validates the process function with deduplication
func TestProcess(t *testing.T) {
	raw := []structs.ApiServer{
		{
			Name:     "us1234",
			Hostname: "us1234.nordvpn.com",
			Station:  "192.168.1.1",
			Load:     50,
			Locations: []structs.Location{
				{
					Country: structs.Country{Name: "US", Code: "US", City: structs.City{Name: "NY"}},
					Lat:     40.0,
					Lon:     -73.0,
				},
			},
			Specs: []structs.Spec{{ID: "version", Values: []structs.SpecVal{{Value: "2.1.0"}}}},
			Tech:  []structs.Tech{{ID: "wireguard_udp", Meta: []structs.Meta{{Name: "public_key", Value: "key1"}}}},
		},
		{
			Name:     "us1234", // Duplicate name
			Hostname: "us1234-2.nordvpn.com",
			Station:  "192.168.1.2",
			Load:     60,
			Locations: []structs.Location{
				{
					Country: structs.Country{Name: "US", Code: "US", City: structs.City{Name: "NY"}},
					Lat:     40.0,
					Lon:     -73.0,
				},
			},
			Specs: []structs.Spec{{ID: "version", Values: []structs.SpecVal{{Value: "2.1.0"}}}},
			Tech:  []structs.Tech{{ID: "wireguard_udp", Meta: []structs.Meta{{Name: "public_key", Value: "key2"}}}},
		},
		{
			Name:     "uk5678",
			Hostname: "uk5678.nordvpn.com",
			Station:  "10.0.0.1",
			Load:     30,
			Locations: []structs.Location{
				{
					Country: structs.Country{Name: "UK", Code: "GB", City: structs.City{Name: "London"}},
					Lat:     51.5,
					Lon:     -0.1,
				},
			},
			Specs: []structs.Spec{{ID: "version", Values: []structs.SpecVal{{Value: "2.2.0"}}}},
			Tech:  []structs.Tech{{ID: "wireguard_udp", Meta: []structs.Meta{{Name: "public_key", Value: "key3"}}}},
		},
	}

	processed, rejected := process(raw, 40.0, -73.0)

	// Should deduplicate us1234 (2 raw servers -> 1 processed)
	if len(processed) != 2 {
		t.Errorf("process() returned %d servers, want 2 (after deduplication)", len(processed))
	}

	if rejected != 1 {
		t.Errorf("process() rejected count = %d, want 1", rejected)
	}

	// Check that both distinct servers are present
	names := make(map[string]bool)
	for _, s := range processed {
		names[s.Name] = true
	}

	if !names["us1234"] || !names["uk5678"] {
		t.Error("process() should contain both us1234 and uk5678 after deduplication")
	}
}
