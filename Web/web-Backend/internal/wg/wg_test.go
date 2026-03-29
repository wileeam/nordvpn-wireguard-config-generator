package wg

import (
	"bytes"
	"strings"
	"testing"

	"nordgen/internal/types"
)

func makeServer(hostname, station string) types.ProcessedServer {
	return types.ProcessedServer{
		Hostname: hostname,
		Station:  station,
	}
}

func makeOpts(key, dns string, useStation bool, keepAlive int) types.ValidatedConfig {
	return types.ValidatedConfig{
		PrivateKey: key,
		DNS:        dns,
		UseStation: useStation,
		KeepAlive:  keepAlive,
	}
}

// ---------------------------------------------------------------------------
// Build
// ---------------------------------------------------------------------------

func TestBuildUsesHostname(t *testing.T) {
	srv := makeServer("us1.nordvpn.com", "203.0.113.5")
	opts := makeOpts("PRIVKEY=", "1.1.1.1", false, 25)
	out := Build(srv, "PUBKEY=", opts)

	if !strings.Contains(string(out), "Endpoint=us1.nordvpn.com:51820") {
		t.Errorf("expected hostname endpoint in:\n%s", out)
	}
}

func TestBuildUsesStation(t *testing.T) {
	srv := makeServer("us1.nordvpn.com", "203.0.113.5")
	opts := makeOpts("PRIVKEY=", "1.1.1.1", true, 25)
	out := Build(srv, "PUBKEY=", opts)

	if !strings.Contains(string(out), "Endpoint=203.0.113.5:51820") {
		t.Errorf("expected station endpoint in:\n%s", out)
	}
}

func TestBuildContainsRequiredSections(t *testing.T) {
	srv := makeServer("us1.nordvpn.com", "1.2.3.4")
	opts := makeOpts("MYKEY=", "8.8.8.8", false, 30)
	out := string(Build(srv, "PUBKEY=", opts))

	checks := []string{
		"[Interface]",
		"PrivateKey=MYKEY=",
		"Address=10.5.0.2/16",
		"DNS=8.8.8.8",
		"[Peer]",
		"PublicKey=PUBKEY=",
		"AllowedIPs=0.0.0.0/0,::/0",
		":51820",
		"PersistentKeepalive=30",
	}
	for _, check := range checks {
		if !strings.Contains(out, check) {
			t.Errorf("Build output missing %q\n\nFull output:\n%s", check, out)
		}
	}
}

func TestBuildKeepaliveValue(t *testing.T) {
	srv := makeServer("host.example.com", "10.0.0.1")
	opts := makeOpts("KEY=", "103.86.96.100", false, 60)
	out := string(Build(srv, "PUB=", opts))

	if !strings.Contains(out, "PersistentKeepalive=60") {
		t.Errorf("expected PersistentKeepalive=60 in:\n%s", out)
	}
}

func TestBuildNonEmpty(t *testing.T) {
	srv := makeServer("host", "1.2.3.4")
	opts := makeOpts("KEY=", "DNS", false, 25)
	out := Build(srv, "PUB=", opts)
	if len(out) == 0 {
		t.Error("Build should return non-empty output")
	}
}

// ---------------------------------------------------------------------------
// WriteConfig
// ---------------------------------------------------------------------------

func TestWriteConfigUsesHostname(t *testing.T) {
	var buf bytes.Buffer
	srv := makeServer("de1.nordvpn.com", "198.51.100.1")
	opts := makeOpts("WK=", "1.0.0.1", false, 25)
	WriteConfig(&buf, srv, "WPUB=", opts)

	if !strings.Contains(buf.String(), "Endpoint=de1.nordvpn.com:51820") {
		t.Errorf("expected hostname endpoint in:\n%s", buf.String())
	}
}

func TestWriteConfigUsesStation(t *testing.T) {
	var buf bytes.Buffer
	srv := makeServer("de1.nordvpn.com", "198.51.100.1")
	opts := makeOpts("WK=", "1.0.0.1", true, 25)
	WriteConfig(&buf, srv, "WPUB=", opts)

	if !strings.Contains(buf.String(), "Endpoint=198.51.100.1:51820") {
		t.Errorf("expected station endpoint in:\n%s", buf.String())
	}
}

func TestWriteConfigMatchesBuild(t *testing.T) {
	srv := makeServer("us5.nordvpn.com", "10.10.10.10")
	opts := makeOpts("KEY=", "8.8.4.4", false, 45)
	pubKey := "PUB="

	var buf bytes.Buffer
	WriteConfig(&buf, srv, pubKey, opts)
	fromWrite := buf.Bytes()

	fromBuild := Build(srv, pubKey, opts)

	if !bytes.Equal(fromWrite, fromBuild) {
		t.Errorf("WriteConfig and Build differ:\nWriteConfig: %q\nBuild: %q", fromWrite, fromBuild)
	}
}

func TestWriteConfigPoolReuse(t *testing.T) {
	// Call WriteConfig multiple times to exercise pool reuse
	srv := makeServer("host", "1.2.3.4")
	opts := makeOpts("KEY=", "DNS", false, 25)

	for i := 0; i < 10; i++ {
		var buf bytes.Buffer
		WriteConfig(&buf, srv, "PUB=", opts)
		if buf.Len() == 0 {
			t.Errorf("iteration %d: WriteConfig produced empty output", i)
		}
	}
}

func TestWriteConfigNonEmpty(t *testing.T) {
	var buf bytes.Buffer
	srv := makeServer("host", "1.2.3.4")
	opts := makeOpts("KEY=", "DNS", false, 25)
	WriteConfig(&buf, srv, "PUB=", opts)
	if buf.Len() == 0 {
		t.Error("WriteConfig should write non-empty output")
	}
}
