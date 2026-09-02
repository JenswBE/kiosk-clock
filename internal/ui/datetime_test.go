package ui

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "00:00:00",
			time: time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC),
			want: "00:00",
		},
		{
			name: "00:00:01",
			time: time.Date(2026, time.March, 2, 0, 0, 1, 0, time.UTC),
			want: "00 00",
		},
		{
			name: "00:00:02",
			time: time.Date(2026, time.March, 2, 0, 0, 2, 0, time.UTC),
			want: "00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTime(tt.time); got != tt.want {
				t.Errorf("formatTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "monday in march",
			time: time.Date(2026, time.March, 2, 12, 0, 0, 0, time.UTC),
			want: "maandag 2 maart 2026",
		},
		{
			name: "sunday in december",
			time: time.Date(2025, time.December, 28, 0, 0, 0, 0, time.UTC),
			want: "zondag 28 december 2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDate(tt.time); got != tt.want {
				t.Errorf("formatDate() = %q, want %q", got, tt.want)
			}
		})
	}
}
