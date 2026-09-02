package system

import "testing"

func TestParseBacklightPercentage(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    uint8
		wantErr bool
	}{
		{
			name:   "single device",
			output: "amdgpu_bl1,backlight,28214,44%,64764\n",
			want:   44,
		},
		{
			name:   "multiple devices uses first non-empty line",
			output: "\n  \nintel_backlight,backlight,100,0%,255\nother,backlight,255,100%,255\n",
			want:   0,
		},
		{
			name:    "empty output",
			output:  "",
			wantErr: true,
		},
		{
			name:    "blank output",
			output:  "   \n  \n",
			wantErr: true,
		},
		{
			name:    "too few fields",
			output:  "amdgpu_bl1,backlight,44%",
			wantErr: true,
		},
		{
			name:    "non numeric percentage",
			output:  "amdgpu_bl1,backlight,28214,abc%,64764",
			wantErr: true,
		},
		{
			name:    "out of range percentage",
			output:  "amdgpu_bl1,backlight,28214,150%,64764",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBacklightPercentage(tt.output)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseBacklightPercentage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && got != tt.want {
				t.Errorf("parseBacklightPercentage() = %d, want %d", got, tt.want)
			}
		})
	}
}
