// Package config handles loading, saving and validating the clock's
// user configuration file.
package config

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const fileName = ".fullscreen-clock.json"

// Config holds the user-configurable settings of the clock.
type Config struct {
	TextColor       string `json:"text_color"`
	BackgroundColor string `json:"background_color"`
}

// Default is the configuration used when no config file exists yet, or
// when it cannot be read/parsed.
var Default = Config{
	TextColor:       "#FFFFFF",
	BackgroundColor: "#000000",
}

// Path returns the location of the config file in the user's home
// directory.
func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return fileName
	}

	return filepath.Join(home, fileName)
}

// Load reads the config file, falling back to Default if it is
// missing, invalid or contains invalid colors.
func Load() Config {
	data, err := os.ReadFile(Path())
	if err != nil {
		return Default
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default
	}

	if _, err := ParseColor(cfg.TextColor); err != nil {
		cfg.TextColor = Default.TextColor
	}

	if _, err := ParseColor(cfg.BackgroundColor); err != nil {
		cfg.BackgroundColor = Default.BackgroundColor
	}

	return cfg
}

// Save writes the config file to disk.
func Save(cfg Config) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(Path(), data, 0o600)
}

// ParseColor parses a "#RRGGBB" hex string into a color.Color.
func ParseColor(value string) (color.Color, error) {
	value = strings.TrimPrefix(value, "#")

	if len(value) != 6 {
		return nil, fmt.Errorf("invalid color: %q", value)
	}

	r, err := strconv.ParseUint(value[0:2], 16, 8)
	if err != nil {
		return nil, err
	}

	g, err := strconv.ParseUint(value[2:4], 16, 8)
	if err != nil {
		return nil, err
	}

	b, err := strconv.ParseUint(value[4:6], 16, 8)
	if err != nil {
		return nil, err
	}

	return color.NRGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}, nil
}

// ColorHex formats a color.Color as a "#RRGGBB" hex string.
func ColorHex(c color.Color) string {
	r, g, b, _ := c.RGBA()

	return fmt.Sprintf(
		"#%02X%02X%02X",
		uint8(r>>8),
		uint8(g>>8),
		uint8(b>>8),
	)
}
