package wulimt

import (
	"testing"
)

func TestParseRate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"1-H", false},
		{"1-h", false},
		{"10-h", false},
		{"10-2h", false},
		{"5-m", false},
		{"5-1m", false},
		{"3-10m", false},
		{"3-12s", false},
		{"3-11.2m", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRate(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("ParseRate() raw = %q, got = %v", tt.name, got)
			}
		})
	}
}
