package config

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    color.NRGBA
		wantErr bool
	}{
		{name: "with hash", value: "#FF0080", want: color.NRGBA{R: 0xFF, G: 0x00, B: 0x80, A: 255}},
		{name: "without hash", value: "00FF00", want: color.NRGBA{R: 0x00, G: 0xFF, B: 0x00, A: 255}},
		{name: "too short", value: "#FFF", wantErr: true},
		{name: "invalid hex", value: "#GGFFFF", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColor(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseColor(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}

			if err == nil && got != tt.want {
				t.Errorf("ParseColor(%q) = %+v, want %+v", tt.value, got, tt.want)
			}
		})
	}
}

func TestColorHex(t *testing.T) {
	got := ColorHex(color.NRGBA{R: 0xFF, G: 0x00, B: 0x80, A: 255})
	if want := "#FF0080"; got != want {
		t.Errorf("ColorHex() = %q, want %q", got, want)
	}
}

func TestLoadSaveRoundTrip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	// No config file yet -> Default.
	if got := Load(); got != Default {
		t.Fatalf("Load() = %+v, want Default %+v", got, Default)
	}

	cfg := Config{TextColor: "#112233", BackgroundColor: "#445566"}
	Save(cfg)

	if got := Load(); got != cfg {
		t.Errorf("Load() = %+v, want %+v", got, cfg)
	}
}

func TestLoadInvalidColorsFallBackToDefault(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	data := []byte(`{"text_color":"not-a-color","background_color":"also-invalid"}`)
	if err := os.WriteFile(filepath.Join(home, fileName), data, 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	got := Load()
	if got.TextColor != Default.TextColor {
		t.Errorf("TextColor = %q, want %q", got.TextColor, Default.TextColor)
	}

	if got.BackgroundColor != Default.BackgroundColor {
		t.Errorf("BackgroundColor = %q, want %q", got.BackgroundColor, Default.BackgroundColor)
	}
}

func TestLoadMissingFileFallsBackToDefault(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if got := Load(); got != Default {
		t.Errorf("Load() = %+v, want Default %+v", got, Default)
	}
}
