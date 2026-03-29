package gen

import (
	"testing"

	"github.com/mustafachyi/nordvpn-wireguard-config-generator/internal/structs"
)

// ---------------------------------------------------------------------------
// sanitize
// ---------------------------------------------------------------------------

func TestSanitize(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"United States", "united_states"},
		{"New York", "new_york"},
		{"", ""},
		{"us", "us"},
		{"US", "us"},
		{"foo<bar>baz", "foobarbaz"},
		{"a:b", "ab"},
		{`a"b`, "ab"},
		{`a/b`, "ab"},
		{`a\b`, "ab"},
		{"a|b", "ab"},
		{"a?b", "ab"},
		{"a*b", "ab"},
		{"a\x00b", "ab"},
		{"a#b", "ab"},
		{"already_lower", "already_lower"},
	}

	for _, tc := range cases {
		got := sanitize(tc.input)
		if got != tc.want {
			t.Errorf("sanitize(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// baseName
// ---------------------------------------------------------------------------

func makeServer(name, code, ip string) structs.Server {
	return structs.Server{Name: name, Code: code, IP: ip}
}

func TestBaseName(t *testing.T) {
	cases := []struct {
		s    structs.Server
		want string
	}{
		{makeServer("United States #1234", "us", "1.2.3.4"), "us1234.conf"},
		{makeServer("US#1", "us", "1.2.3.4"), "us1.conf"},
		{makeServer("DE#999", "de", "5.6.7.8"), "de999.conf"},
	}

	for _, tc := range cases {
		got := baseName(tc.s)
		if got != tc.want {
			t.Errorf("baseName(%q) = %q, want %q", tc.s.Name, got, tc.want)
		}
	}
}

func TestBaseNameEndsWithConf(t *testing.T) {
	s := makeServer("DE #99", "de", "1.2.3.4")
	result := baseName(s)
	if len(result) < 5 || result[len(result)-5:] != ".conf" {
		t.Errorf("baseName should end with .conf, got %q", result)
	}
}

func TestBaseNameNoDigitsFallback(t *testing.T) {
	s := makeServer("NoNumbers", "us", "192.168.1.1")
	result := baseName(s)
	if len(result) < 5 || result[len(result)-5:] != ".conf" {
		t.Errorf("fallback should end with .conf, got %q", result)
	}
	// base without .conf should not contain dots
	base := result[:len(result)-5]
	for _, c := range base {
		if c == '.' {
			t.Errorf("baseName fallback should strip dots from IP, got %q", result)
		}
	}
}

func TestBaseNameMaxLength(t *testing.T) {
	// Create a name with many digits
	s := makeServer("Server #123456789012345", "us", "1.2.3.4")
	result := baseName(s)
	base := result[:len(result)-5] // strip .conf
	if len(base) > 15 {
		t.Errorf("baseName base exceeds 15 chars: %q (len=%d)", base, len(base))
	}
}

func TestBaseNameFallbackMaxLength(t *testing.T) {
	s := makeServer("NoNumbers", "us", "192.168.100.200")
	result := baseName(s)
	base := result[:len(result)-5]
	if len(base) > 15 {
		t.Errorf("fallback base exceeds 15 chars: %q (len=%d)", base, len(base))
	}
}

// ---------------------------------------------------------------------------
// buildConfig (Writer)
// ---------------------------------------------------------------------------

func makeTestWriter(useIP bool) *Writer {
	return &Writer{
		key: "MYPRIVKEY=",
		prefs: structs.Preferences{
			DNS:       "1.1.1.1",
			UseIP:     useIP,
			Keepalive: 25,
		},
	}
}

func makeTestServer() structs.Server {
	return structs.Server{
		Name:    "US#1",
		Host:    "us1.nordvpn.com",
		IP:      "203.0.113.5",
		Load:    20,
		Country: "United States",
		Code:    "us",
		City:    "New York",
		PubK:    "PUBKEY123=",
	}
}

func TestBuildConfigUsesHostname(t *testing.T) {
	w := makeTestWriter(false)
	s := makeTestServer()
	cfg := w.buildConfig(s)

	want := "[Interface]\nPrivateKey = MYPRIVKEY=\nAddress = 10.5.0.2/16\nDNS = 1.1.1.1\n\n" +
		"[Peer]\nPublicKey = PUBKEY123=\nAllowedIPs = 0.0.0.0/0, ::/0\n" +
		"Endpoint = us1.nordvpn.com:51820\nPersistentKeepalive = 25"
	if cfg != want {
		t.Errorf("buildConfig (hostname):\ngot:  %q\nwant: %q", cfg, want)
	}
}

func TestBuildConfigUsesIP(t *testing.T) {
	w := makeTestWriter(true)
	s := makeTestServer()
	cfg := w.buildConfig(s)

	want := "[Interface]\nPrivateKey = MYPRIVKEY=\nAddress = 10.5.0.2/16\nDNS = 1.1.1.1\n\n" +
		"[Peer]\nPublicKey = PUBKEY123=\nAllowedIPs = 0.0.0.0/0, ::/0\n" +
		"Endpoint = 203.0.113.5:51820\nPersistentKeepalive = 25"
	if cfg != want {
		t.Errorf("buildConfig (IP):\ngot:  %q\nwant: %q", cfg, want)
	}
}

func TestBuildConfigContainsRequiredSections(t *testing.T) {
	w := makeTestWriter(false)
	s := makeTestServer()
	cfg := w.buildConfig(s)

	checks := []string{
		"[Interface]",
		"[Peer]",
		"PrivateKey = MYPRIVKEY=",
		"Address = 10.5.0.2/16",
		"DNS = 1.1.1.1",
		"PublicKey = PUBKEY123=",
		"AllowedIPs = 0.0.0.0/0, ::/0",
		":51820",
		"PersistentKeepalive = 25",
	}
	for _, check := range checks {
		found := false
		for i := 0; i <= len(cfg)-len(check); i++ {
			if cfg[i:i+len(check)] == check {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("buildConfig missing %q in output:\n%s", check, cfg)
		}
	}
}

func TestBuildConfigDifferentKeepalive(t *testing.T) {
	w := &Writer{
		key: "KEY=",
		prefs: structs.Preferences{
			DNS:       "8.8.8.8",
			UseIP:     false,
			Keepalive: 60,
		},
	}
	s := makeTestServer()
	cfg := w.buildConfig(s)

	keepaliveStr := "PersistentKeepalive = 60"
	found := false
	for i := 0; i <= len(cfg)-len(keepaliveStr); i++ {
		if cfg[i:i+len(keepaliveStr)] == keepaliveStr {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected %q in config, got:\n%s", keepaliveStr, cfg)
	}
}
