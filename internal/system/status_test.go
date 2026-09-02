package system

import "testing"

func TestParseNTPSynchronized(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   bool
	}{
		{name: "synced", output: "yes\n", want: true},
		{name: "not synced", output: "no\n", want: false},
		{name: "empty", output: "", want: false},
		{name: "padded", output: "  yes  \n", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNTPSynchronized(tt.output); got != tt.want {
				t.Errorf("parseNTPSynchronized(%q) = %v, want %v", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseWifiConnected(t *testing.T) {
	tests := []struct {
		name       string
		output     string
		wantStatus string
		wantOK     bool
	}{
		{
			name:       "wifi connected",
			output:     "ethernet:unavailable\nwifi:connected\nloopback:unmanaged\n",
			wantStatus: "connected",
			wantOK:     true,
		},
		{
			name:       "wifi disconnected",
			output:     "ethernet:unavailable\nwifi:disconnected\n",
			wantStatus: "disconnected",
			wantOK:     false,
		},
		{
			name:       "no wifi device",
			output:     "ethernet:connected\n",
			wantStatus: "",
			wantOK:     false,
		},
		{
			name:       "empty output",
			output:     "",
			wantStatus: "",
			wantOK:     false,
		},
		{
			name:       "malformed line ignored",
			output:     "malformed-line\nwifi:connected\n",
			wantStatus: "connected",
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStatus, gotOK := parseWifiConnected(tt.output); gotStatus != tt.wantStatus || gotOK != tt.wantOK {
				t.Errorf("parseWifiConnected(%q) = (%v, %v), want (%v, %v)", tt.output, gotStatus, gotOK, tt.wantStatus, tt.wantOK)
			}
		})
	}
}
