package wg

import (
	"bytes"
	"strings"
	"testing"

	"nordgen/internal/types"
)

// TestBuild validates the WireGuard configuration building function
func TestBuild(t *testing.T) {
	tests := []struct {
		name       string
		server     types.ProcessedServer
		pubKey     string
		opts       types.ValidatedConfig
		wantFields map[string]string
	}{
		{
			name: "basic config with hostname",
			server: types.ProcessedServer{
				Hostname: "us1234.nordvpn.com",
				Station:  "192.168.1.1",
			},
			pubKey: "serverPublicKey123456789012345678901234=",
			opts: types.ValidatedConfig{
				PrivateKey: "userPrivateKey1234567890123456789012345=",
				DNS:        "103.86.96.100",
				UseStation: false,
				KeepAlive:  25,
			},
			wantFields: map[string]string{
				"PrivateKey":           "userPrivateKey1234567890123456789012345=",
				"Address":              "10.5.0.2/16",
				"DNS":                  "103.86.96.100",
				"PublicKey":            "serverPublicKey123456789012345678901234=",
				"AllowedIPs":           "0.0.0.0/0,::/0",
				"Endpoint":             "us1234.nordvpn.com:51820",
				"PersistentKeepalive": "25",
			},
		},
		{
			name: "config with station IP",
			server: types.ProcessedServer{
				Hostname: "us1234.nordvpn.com",
				Station:  "192.168.1.1",
			},
			pubKey: "serverPublicKey123456789012345678901234=",
			opts: types.ValidatedConfig{
				PrivateKey: "userPrivateKey1234567890123456789012345=",
				DNS:        "1.1.1.1,8.8.8.8",
				UseStation: true,
				KeepAlive:  60,
			},
			wantFields: map[string]string{
				"PrivateKey":           "userPrivateKey1234567890123456789012345=",
				"Address":              "10.5.0.2/16",
				"DNS":                  "1.1.1.1,8.8.8.8",
				"PublicKey":            "serverPublicKey123456789012345678901234=",
				"AllowedIPs":           "0.0.0.0/0,::/0",
				"Endpoint":             "192.168.1.1:51820",
				"PersistentKeepalive": "60",
			},
		},
		{
			name: "config with custom DNS",
			server: types.ProcessedServer{
				Hostname: "uk5678.nordvpn.com",
				Station:  "10.0.0.1",
			},
			pubKey: "anotherPublicKey234567890123456789012345=",
			opts: types.ValidatedConfig{
				PrivateKey: "anotherPrivateKey23456789012345678901234=",
				DNS:        "9.9.9.9",
				UseStation: false,
				KeepAlive:  15,
			},
			wantFields: map[string]string{
				"PrivateKey":           "anotherPrivateKey23456789012345678901234=",
				"Address":              "10.5.0.2/16",
				"DNS":                  "9.9.9.9",
				"PublicKey":            "anotherPublicKey234567890123456789012345=",
				"AllowedIPs":           "0.0.0.0/0,::/0",
				"Endpoint":             "uk5678.nordvpn.com:51820",
				"PersistentKeepalive": "15",
			},
		},
		{
			name: "config with max keepalive",
			server: types.ProcessedServer{
				Hostname: "jp9999.nordvpn.com",
				Station:  "172.16.0.1",
			},
			pubKey: "maxKeepalivePublicKey123456789012345678=",
			opts: types.ValidatedConfig{
				PrivateKey: "maxKeepalivePrivateKey12345678901234567=",
				DNS:        "8.8.4.4",
				UseStation: false,
				KeepAlive:  120,
			},
			wantFields: map[string]string{
				"PersistentKeepalive": "120",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Build(tt.server, tt.pubKey, tt.opts)
			resultStr := string(result)

			// Check for [Interface] and [Peer] sections
			if !strings.Contains(resultStr, "[Interface]") {
				t.Error("Config missing [Interface] section")
			}
			if !strings.Contains(resultStr, "[Peer]") {
				t.Error("Config missing [Peer] section")
			}

			// Check all expected fields
			for field, value := range tt.wantFields {
				expectedLine := field + "=" + value
				if !strings.Contains(resultStr, expectedLine) {
					t.Errorf("Config missing or incorrect field %q: expected %q\nGot config:\n%s", field, expectedLine, resultStr)
				}
			}

			// Verify structure
			if !strings.HasPrefix(resultStr, "[Interface]\n") {
				t.Error("Config should start with [Interface]")
			}
		})
	}
}

// TestWriteConfig validates the WriteConfig function that writes to an io.Writer
func TestWriteConfig(t *testing.T) {
	server := types.ProcessedServer{
		Hostname: "test.nordvpn.com",
		Station:  "192.168.100.1",
	}
	pubKey := "testPublicKey3456789012345678901234567890=",
	opts := types.ValidatedConfig{
		PrivateKey: "testPrivateKey456789012345678901234567=",
		DNS:        "1.1.1.1",
		UseStation: false,
		KeepAlive:  30,
	}

	var buf bytes.Buffer
	WriteConfig(&buf, server, pubKey, opts)

	result := buf.String()

	// Basic structure checks
	if !strings.Contains(result, "[Interface]") {
		t.Error("Config missing [Interface] section")
	}
	if !strings.Contains(result, "[Peer]") {
		t.Error("Config missing [Peer] section")
	}
	if !strings.Contains(result, "PrivateKey=testPrivateKey456789012345678901234567=") {
		t.Error("Config missing correct PrivateKey")
	}
	if !strings.Contains(result, "PublicKey=testPublicKey3456789012345678901234567890=") {
		t.Error("Config missing correct PublicKey")
	}
	if !strings.Contains(result, "DNS=1.1.1.1") {
		t.Error("Config missing correct DNS")
	}
	if !strings.Contains(result, "Endpoint=test.nordvpn.com:51820") {
		t.Error("Config missing correct Endpoint")
	}
	if !strings.Contains(result, "PersistentKeepalive=30") {
		t.Error("Config missing correct PersistentKeepalive")
	}
}

// TestWriteConfigWithStation validates WriteConfig when using station IP
func TestWriteConfigWithStation(t *testing.T) {
	server := types.ProcessedServer{
		Hostname: "test.nordvpn.com",
		Station:  "10.20.30.40",
	}
	pubKey := "stationPublicKey456789012345678901234567=",
	opts := types.ValidatedConfig{
		PrivateKey: "stationPrivateKey56789012345678901234567=",
		DNS:        "8.8.8.8",
		UseStation: true,
		KeepAlive:  45,
	}

	var buf bytes.Buffer
	WriteConfig(&buf, server, pubKey, opts)

	result := buf.String()

	if !strings.Contains(result, "Endpoint=10.20.30.40:51820") {
		t.Errorf("Config should use station IP as endpoint, got:\n%s", result)
	}
}

// TestBuildAndWriteEquivalence ensures Build and WriteConfig produce the same output
func TestBuildAndWriteEquivalence(t *testing.T) {
	server := types.ProcessedServer{
		Hostname: "equiv.nordvpn.com",
		Station:  "172.31.0.1",
	}
	pubKey := "equivPublicKey6789012345678901234567890=",
	opts := types.ValidatedConfig{
		PrivateKey: "equivPrivateKey789012345678901234567890=",
		DNS:        "9.9.9.9,1.1.1.1",
		UseStation: false,
		KeepAlive:  50,
	}

	// Test Build
	buildResult := Build(server, pubKey, opts)

	// Test WriteConfig
	var buf bytes.Buffer
	WriteConfig(&buf, server, pubKey, opts)
	writeResult := buf.Bytes()

	// They should produce identical output
	if !bytes.Equal(buildResult, writeResult) {
		t.Errorf("Build and WriteConfig produced different outputs:\nBuild:\n%s\n\nWriteConfig:\n%s", buildResult, writeResult)
	}
}
