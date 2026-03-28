package main

import (
	"testing"
)

// TestIsHex validates the hex string validation function
func TestIsHex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid 64 char hex lowercase", "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789", true},
		{"valid 64 char hex uppercase", "ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789", true},
		{"valid 64 char hex mixed", "ABcdEF0123456789abCDef0123456789ABcdEF0123456789abCDef0123456789", true},
		{"too short", "abcdef01234567", false},
		{"too long", "abcdef0123456789abcdef0123456789abcdef0123456789abcdef01234567890", false},
		{"invalid char g", "gbcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789", false},
		{"invalid char z", "abcdef0123456789abcdef0123456789abcdef0123456789abcdef012345678z", false},
		{"empty string", "", false},
		{"63 chars", "abcdef0123456789abcdef0123456789abcdef0123456789abcdef012345678", false},
		{"65 chars", "abcdef0123456789abcdef0123456789abcdef0123456789abcdef01234567890", false},
		{"special characters", "abcdef0123456789-bcdef0123456789abcdef0123456789abcdef0123456789", false},
		{"spaces", "abcdef0123456789 bcdef0123456789abcdef0123456789abcdef0123456789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isHex(tt.input)
			if got != tt.want {
				t.Errorf("isHex(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestIsKey validates the WireGuard key validation function
func TestIsKey(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid key", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=", true},
		{"valid key with numbers", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/abcde=", true},
		{"valid key with plus", "ABCDEFGHIJKLMNOPQRSTUVWXYZ++++++++++abcdefg=", true},
		{"valid key with slash", "ABCDEFGHIJKLMNOPQRSTUVWXYZ//////////abcdefg=", true},
		{"too short", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmn=", false},
		{"too long", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnop=", false},
		{"missing equals", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnop", false},
		{"equals not at end", "=BCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=", false},
		{"invalid char", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmn@=", false},
		{"empty string", "", false},
		{"43 chars no equals", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnop", false},
		{"double equals", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno==", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isKey(tt.input)
			if got != tt.want {
				t.Errorf("isKey(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestIsIPv4 validates the IPv4 address validation function
func TestIsIPv4(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid IP", "192.168.1.1", true},
		{"valid IP zeros", "0.0.0.0", true},
		{"valid IP max", "255.255.255.255", true},
		{"valid IP mixed", "10.20.30.40", true},
		{"valid IP single digits", "1.2.3.4", true},
		{"empty string", "", false},
		{"too few octets", "192.168.1", false},
		{"too many octets", "192.168.1.1.1", false},
		{"octet out of range 256", "192.168.1.256", false},
		{"octet out of range 300", "300.168.1.1", false},
		{"negative number", "192.168.-1.1", false},
		{"letters", "192.168.a.1", false},
		{"spaces", "192.168. 1.1", false},
		{"leading zero valid", "192.168.001.1", false}, // Custom validator considers this invalid due to octet > 255 during parsing
		{"trailing dot", "192.168.1.1.", false},
		{"leading dot", ".192.168.1.1", false},
		{"double dot", "192..168.1.1", false},
		{"no dots", "192168", false},
		{"only dots", "...", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIPv4(tt.input)
			if got != tt.want {
				t.Errorf("isIPv4(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestParseCommon validates the common configuration parsing function
func TestParseCommon(t *testing.T) {
	keepAlive20 := 20
	keepAlive10 := 10
	keepAlive150 := 150

	tests := []struct {
		name          string
		key           string
		dns           string
		endpoint      string
		keepAlive     *int
		wantErrs      int
		wantDNS       string
		wantStation   bool
		wantKeepAlive int
	}{
		{
			name:          "all valid",
			key:           "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=",
			dns:           "1.1.1.1,8.8.8.8",
			endpoint:      "hostname",
			keepAlive:     &keepAlive20,
			wantErrs:      0,
			wantDNS:       "1.1.1.1,8.8.8.8",
			wantStation:   false,
			wantKeepAlive: 20,
		},
		{
			name:          "default values",
			key:           "",
			dns:           "",
			endpoint:      "",
			keepAlive:     nil,
			wantErrs:      0,
			wantDNS:       "103.86.96.100",
			wantStation:   false,
			wantKeepAlive: 25,
		},
		{
			name:          "station endpoint",
			key:           "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=",
			dns:           "1.1.1.1",
			endpoint:      "station",
			keepAlive:     &keepAlive20,
			wantErrs:      0,
			wantDNS:       "1.1.1.1",
			wantStation:   true,
			wantKeepAlive: 20,
		},
		{
			name:          "invalid key",
			key:           "invalid_key",
			dns:           "1.1.1.1",
			endpoint:      "hostname",
			keepAlive:     &keepAlive20,
			wantErrs:      1,
			wantDNS:       "1.1.1.1",
			wantStation:   false,
			wantKeepAlive: 20,
		},
		{
			name:          "invalid DNS",
			key:           "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=",
			dns:           "256.1.1.1",
			endpoint:      "hostname",
			keepAlive:     &keepAlive20,
			wantErrs:      1,
			wantDNS:       "103.86.96.100",
			wantStation:   false,
			wantKeepAlive: 20,
		},
		{
			name:          "invalid endpoint",
			key:           "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=",
			dns:           "1.1.1.1",
			endpoint:      "invalid",
			keepAlive:     &keepAlive20,
			wantErrs:      1,
			wantDNS:       "1.1.1.1",
			wantStation:   false,
			wantKeepAlive: 20,
		},
		{
			name:          "keepalive too low",
			key:           "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=",
			dns:           "1.1.1.1",
			endpoint:      "hostname",
			keepAlive:     &keepAlive10,
			wantErrs:      1,
			wantDNS:       "1.1.1.1",
			wantStation:   false,
			wantKeepAlive: 25,
		},
		{
			name:          "keepalive too high",
			key:           "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=",
			dns:           "1.1.1.1",
			endpoint:      "hostname",
			keepAlive:     &keepAlive150,
			wantErrs:      1,
			wantDNS:       "1.1.1.1",
			wantStation:   false,
			wantKeepAlive: 25,
		},
		{
			name:          "multiple invalid fields",
			key:           "invalid",
			dns:           "999.999.999.999",
			endpoint:      "bad",
			keepAlive:     &keepAlive10,
			wantErrs:      4,
			wantDNS:       "103.86.96.100",
			wantStation:   false,
			wantKeepAlive: 25,
		},
		{
			name:          "DNS with spaces",
			key:           "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno=",
			dns:           " 1.1.1.1 , 8.8.8.8 ",
			endpoint:      "hostname",
			keepAlive:     &keepAlive20,
			wantErrs:      0,
			wantDNS:       " 1.1.1.1 , 8.8.8.8 ",
			wantStation:   false,
			wantKeepAlive: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, errs := parseCommon(tt.key, tt.dns, tt.endpoint, tt.keepAlive)

			if len(errs) != tt.wantErrs {
				t.Errorf("parseCommon() error count = %d, want %d (errors: %v)", len(errs), tt.wantErrs, errs)
			}

			if cfg.DNS != tt.wantDNS {
				t.Errorf("parseCommon() DNS = %q, want %q", cfg.DNS, tt.wantDNS)
			}

			if cfg.UseStation != tt.wantStation {
				t.Errorf("parseCommon() UseStation = %v, want %v", cfg.UseStation, tt.wantStation)
			}

			if cfg.KeepAlive != tt.wantKeepAlive {
				t.Errorf("parseCommon() KeepAlive = %d, want %d", cfg.KeepAlive, tt.wantKeepAlive)
			}
		})
	}
}

// TestExtractRefererHost validates the referer host extraction function
func TestExtractRefererHost(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"http with path", "http://example.com/path", "example.com"},
		{"https with path", "https://example.com/path/to/resource", "example.com"},
		{"http with port", "http://example.com:8080/path", "example.com"},
		{"https with port", "https://example.com:443/path", "example.com"},
		{"no scheme", "example.com/path", "example.com"},
		{"no path", "https://example.com", "example.com"},
		{"subdomain", "https://api.example.com/v1/resource", "api.example.com"},
		{"with query", "https://example.com/path?query=value", "example.com"},
		{"with fragment", "https://example.com/path#fragment", "example.com"},
		{"localhost", "http://localhost:3000/", "localhost"},
		{"IP address", "http://192.168.1.1:8080/api", "192.168.1.1"},
		{"empty string", "", ""},
		{"just domain", "example.com", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRefererHost(tt.input)
			if got != tt.want {
				t.Errorf("extractRefererHost(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestSanitizeFilename validates the filename sanitization function
func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"alphanumeric", "abc123", "abc123"},
		{"uppercase", "ABC", "ABC"},
		{"with spaces", "hello world", "hello_world"},
		{"with dashes", "hello-world", "hello-world"},
		{"with underscores", "hello_world", "hello_world"},
		{"special chars", "hello@#$%world", "helloworld"},
		{"multiple spaces", "hello   world", "hello___world"},
		{"mixed valid and invalid", "Test-File_123.txt", "Test-File_123txt"},
		{"only special chars", "@#$%^&*()", ""},
		{"empty string", "", ""},
		{"unicode", "café", "caf"},
		{"path separators", "path/to/file", "pathtofile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestBuildBaseName validates the base name building function
func TestBuildBaseName(t *testing.T) {
	tests := []struct {
		name    string
		country string
		city    string
		want    string
	}{
		{"no filters", "", "", "NordVPN_All"},
		{"country only", "United States", "", "NordVPN_United_States"},
		{"country and city", "United States", "New York", "NordVPN_United_States_New_York"},
		{"country with spaces", "New Zealand", "", "NordVPN_New_Zealand"},
		{"city with special chars", "United States", "Los Angeles", "NordVPN_United_States_Los_Angeles"},
		{"single word country", "Japan", "", "NordVPN_Japan"},
		{"single word city", "Japan", "Tokyo", "NordVPN_Japan_Tokyo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildBaseName(tt.country, tt.city)
			if got != tt.want {
				t.Errorf("buildBaseName(%q, %q) = %q, want %q", tt.country, tt.city, got, tt.want)
			}
		})
	}
}

// TestBuildDisposition validates the Content-Disposition header building
func TestBuildDisposition(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple name", "test", `attachment; filename="test.nord"`},
		{"with underscore", "test_file", `attachment; filename="test_file.nord"`},
		{"empty name", "", `attachment; filename=".nord"`},
		{"long name", "verylongfilenamehere", `attachment; filename="verylongfilenamehere.nord"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDisposition(tt.input)
			if got != tt.want {
				t.Errorf("buildDisposition(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestBuildConfDisposition validates the WireGuard config Content-Disposition header building
func TestBuildConfDisposition(t *testing.T) {
	tests := []struct {
		name string
		code string
		num  string
		want string
	}{
		{"standard", "us", "1234", `attachment; filename="us1234.conf"`},
		{"empty code", "", "1234", `attachment; filename="1234.conf"`},
		{"empty num", "us", "", `attachment; filename="us.conf"`},
		{"both empty", "", "", `attachment; filename=".conf"`},
		{"long values", "abc", "9999", `attachment; filename="abc9999.conf"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildConfDisposition(tt.code, tt.num)
			if got != tt.want {
				t.Errorf("buildConfDisposition(%q, %q) = %q, want %q", tt.code, tt.num, got, tt.want)
			}
		})
	}
}

// TestDedup validates the file path deduplication function
func TestDedup(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		usedPaths map[string]int
		want      string
	}{
		{
			name:      "first occurrence",
			path:      "test.conf",
			usedPaths: map[string]int{},
			want:      "test.conf",
		},
		{
			name:      "second occurrence",
			path:      "test.conf",
			usedPaths: map[string]int{"test.conf": 0},
			want:      "test_1.conf",
		},
		{
			name:      "third occurrence",
			path:      "test.conf",
			usedPaths: map[string]int{"test.conf": 2, "test_1.conf": 0},
			want:      "test_2.conf",
		},
		{
			name:      "gap in sequence",
			path:      "test.conf",
			usedPaths: map[string]int{"test.conf": 3, "test_1.conf": 0, "test_3.conf": 0},
			want:      "test_2.conf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dedup(tt.path, tt.usedPaths)
			if got != tt.want {
				t.Errorf("dedup(%q, %v) = %q, want %q", tt.path, tt.usedPaths, got, tt.want)
			}
		})
	}
}
