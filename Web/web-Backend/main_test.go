package main

import (
	"testing"

	"nordgen/internal/types"
)

// ---------------------------------------------------------------------------
// isHex
// ---------------------------------------------------------------------------

func TestIsHex(t *testing.T) {
	validLower := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	validUpper := "27810D1FE92A3DA65DCE7755E567B25F72B4F61AE434AF050815C5F73FBB4FC1"
	if len(validLower) != 64 {
		t.Fatalf("validLower must be 64 chars, got %d", len(validLower))
	}
	if len(validUpper) != 64 {
		t.Fatalf("validUpper must be 64 chars, got %d", len(validUpper))
	}

	cases := []struct {
		s    string
		want bool
	}{
		{validLower, true},
		{validUpper, true},
		{"", false},
		{validLower[:63], false},
		{validLower + "a", false},
		{"g" + validLower[1:], false},
		{"z" + validLower[1:], false},
		{"000000000000000000000000000000000000000000000000000000000000000g", false},
	}

	for _, tc := range cases {
		got := isHex(tc.s)
		if got != tc.want {
			t.Errorf("isHex(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// isKey
// ---------------------------------------------------------------------------

func TestIsKey(t *testing.T) {
	// A valid WireGuard key is 44 chars: 43 chars of [A-Za-z0-9+/] followed by '='
	validKey := "QGWezOjmUUJkSGt4TqFWqHaGs5v0a2Atg050X+Ya6Vw="
	if len(validKey) != 44 {
		t.Fatalf("test key must be 44 chars, got %d", len(validKey))
	}

	// Another valid key with only alphanumeric chars before =
	alphaNumKey := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefg="
	if len(alphaNumKey) != 44 {
		t.Fatalf("alphaNumKey must be 44 chars, got %d", len(alphaNumKey))
	}

	cases := []struct {
		s    string
		want bool
	}{
		{validKey, true},
		{alphaNumKey, true},
		{"", false},
		{validKey[:43], false},  // too short, no trailing =
		{validKey + "=", false}, // too long (45 chars)
		{"!BCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=", false}, // invalid char at pos 0
		{validKey[:43] + "x", false}, // doesn't end with =
	}

	for _, tc := range cases {
		got := isKey(tc.s)
		if got != tc.want {
			t.Errorf("isKey(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// isIPv4
// ---------------------------------------------------------------------------

func TestIsIPv4(t *testing.T) {
	cases := []struct {
		s    string
		want bool
	}{
		{"192.168.1.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"103.86.96.100", true},
		{"8.8.8.8", true},
		{"", false},
		{"256.0.0.1", false},
		{"192.168.1", false},
		{"192.168.1.1.1", false},
		{"abc.def.ghi.jkl", false},
		{"1.2.3.", false},
		{".1.2.3", false},
		{"1.2..3", false},
		{"1.2.3.256", false},
		{"999.999.999.999", false},
		{"1.2.3.-1", false},
	}

	for _, tc := range cases {
		got := isIPv4(tc.s)
		if got != tc.want {
			t.Errorf("isIPv4(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// parseCommon
// ---------------------------------------------------------------------------

func intPtr(v int) *int { return &v }

func TestParseCommonDefaults(t *testing.T) {
	cfg, errs := parseCommon("", "", "", nil)
	if len(errs) != 0 {
		t.Errorf("expected no errors for empty inputs, got %v", errs)
	}
	if cfg.DNS != "103.86.96.100" {
		t.Errorf("expected default DNS, got %q", cfg.DNS)
	}
	if cfg.KeepAlive != 25 {
		t.Errorf("expected default keepalive=25, got %d", cfg.KeepAlive)
	}
	if cfg.UseStation {
		t.Error("expected UseStation=false by default")
	}
}

func TestParseCommonValidKey(t *testing.T) {
	key := "QGWezOjmUUJkSGt4TqFWqHaGs5v0a2Atg050X+Ya6Vw="
	cfg, errs := parseCommon(key, "", "", nil)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid key, got %v", errs)
	}
	if cfg.PrivateKey != key {
		t.Errorf("PrivateKey not set correctly")
	}
}

func TestParseCommonInvalidKey(t *testing.T) {
	_, errs := parseCommon("notavalidkey", "", "", nil)
	if len(errs) == 0 {
		t.Error("expected error for invalid key")
	}
}

func TestParseCommonValidDNS(t *testing.T) {
	cfg, errs := parseCommon("", "8.8.8.8", "", nil)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid DNS, got %v", errs)
	}
	if cfg.DNS != "8.8.8.8" {
		t.Errorf("DNS not set correctly: %q", cfg.DNS)
	}
}

func TestParseCommonMultiDNS(t *testing.T) {
	cfg, errs := parseCommon("", "8.8.8.8,1.1.1.1", "", nil)
	if len(errs) != 0 {
		t.Errorf("expected no errors for multi DNS, got %v", errs)
	}
	if cfg.DNS != "8.8.8.8,1.1.1.1" {
		t.Errorf("DNS not set correctly: %q", cfg.DNS)
	}
}

func TestParseCommonInvalidDNS(t *testing.T) {
	_, errs := parseCommon("", "notanip", "", nil)
	if len(errs) == 0 {
		t.Error("expected error for invalid DNS")
	}
}

func TestParseCommonEndpointHostname(t *testing.T) {
	cfg, errs := parseCommon("", "", "hostname", nil)
	if len(errs) != 0 {
		t.Errorf("expected no errors for 'hostname' endpoint, got %v", errs)
	}
	if cfg.UseStation {
		t.Error("UseStation should be false for 'hostname'")
	}
}

func TestParseCommonEndpointStation(t *testing.T) {
	cfg, errs := parseCommon("", "", "station", nil)
	if len(errs) != 0 {
		t.Errorf("expected no errors for 'station' endpoint, got %v", errs)
	}
	if !cfg.UseStation {
		t.Error("UseStation should be true for 'station'")
	}
}

func TestParseCommonInvalidEndpoint(t *testing.T) {
	_, errs := parseCommon("", "", "invalid", nil)
	if len(errs) == 0 {
		t.Error("expected error for invalid endpoint type")
	}
}

func TestParseCommonValidKeepalive(t *testing.T) {
	cfg, errs := parseCommon("", "", "", intPtr(60))
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
	if cfg.KeepAlive != 60 {
		t.Errorf("expected keepalive=60, got %d", cfg.KeepAlive)
	}
}

func TestParseCommonKeepaliveBoundaries(t *testing.T) {
	cfg15, errs15 := parseCommon("", "", "", intPtr(15))
	if len(errs15) != 0 {
		t.Errorf("keepalive=15 should be valid, got %v", errs15)
	}
	if cfg15.KeepAlive != 15 {
		t.Errorf("expected 15, got %d", cfg15.KeepAlive)
	}

	cfg120, errs120 := parseCommon("", "", "", intPtr(120))
	if len(errs120) != 0 {
		t.Errorf("keepalive=120 should be valid, got %v", errs120)
	}
	if cfg120.KeepAlive != 120 {
		t.Errorf("expected 120, got %d", cfg120.KeepAlive)
	}
}

func TestParseCommonKeepaliveOutOfRange(t *testing.T) {
	_, errs14 := parseCommon("", "", "", intPtr(14))
	if len(errs14) == 0 {
		t.Error("expected error for keepalive=14")
	}

	_, errs121 := parseCommon("", "", "", intPtr(121))
	if len(errs121) == 0 {
		t.Error("expected error for keepalive=121")
	}
}

func TestParseCommonMultipleErrors(t *testing.T) {
	_, errs := parseCommon("badkey", "badip", "badendpoint", intPtr(5))
	if len(errs) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(errs), errs)
	}
}

// ---------------------------------------------------------------------------
// validateConfig
// ---------------------------------------------------------------------------

func TestValidateConfigMissingFields(t *testing.T) {
	req := types.ConfigRequest{}
	_, errStr := validateConfig(req)
	if errStr == "" {
		t.Error("expected errors for empty ConfigRequest")
	}
}

func TestValidateConfigValid(t *testing.T) {
	req := types.ConfigRequest{
		Country: "united_states",
		City:    "new_york",
		Name:    "US#1",
	}
	cfg, errStr := validateConfig(req)
	if errStr != "" {
		t.Errorf("unexpected error: %s", errStr)
	}
	if cfg.Name != "US#1" {
		t.Errorf("Name not set: %q", cfg.Name)
	}
}

func TestValidateConfigMissingCountry(t *testing.T) {
	req := types.ConfigRequest{City: "new_york", Name: "US#1"}
	_, errStr := validateConfig(req)
	if errStr == "" {
		t.Error("expected error for missing country")
	}
}

func TestValidateConfigMissingCity(t *testing.T) {
	req := types.ConfigRequest{Country: "united_states", Name: "US#1"}
	_, errStr := validateConfig(req)
	if errStr == "" {
		t.Error("expected error for missing city")
	}
}

func TestValidateConfigMissingName(t *testing.T) {
	req := types.ConfigRequest{Country: "united_states", City: "new_york"}
	_, errStr := validateConfig(req)
	if errStr == "" {
		t.Error("expected error for missing name")
	}
}

// ---------------------------------------------------------------------------
// extractRefererHost
// ---------------------------------------------------------------------------

func TestExtractRefererHost(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://example.com/path", "example.com"},
		{"http://example.com/path", "example.com"},
		{"https://example.com:8080/path", "example.com"},
		{"example.com/path", "example.com"},
		{"https://sub.example.com/a/b/c", "sub.example.com"},
	}

	for _, tc := range cases {
		got := extractRefererHost(tc.in)
		if got != tc.want {
			t.Errorf("extractRefererHost(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// buildBatchPath
// ---------------------------------------------------------------------------

func TestBuildBatchPath(t *testing.T) {
	srv := types.ProcessedServer{
		Country:  "united_states",
		City:     "new_york",
		FileName: "us1234.conf",
	}

	t.Run("no_filter", func(t *testing.T) {
		p := buildBatchPath("", "", srv)
		want := "united_states/new_york/us1234.conf"
		if p != want {
			t.Errorf("got %q, want %q", p, want)
		}
	})

	t.Run("country_only", func(t *testing.T) {
		p := buildBatchPath("united_states", "", srv)
		want := "new_york/us1234.conf"
		if p != want {
			t.Errorf("got %q, want %q", p, want)
		}
	})

	t.Run("country_and_city", func(t *testing.T) {
		p := buildBatchPath("united_states", "new_york", srv)
		want := "us1234.conf"
		if p != want {
			t.Errorf("got %q, want %q", p, want)
		}
	})
}
