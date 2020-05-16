package wulimt

import (
	"testing"
)

func TestRateLimiter_String(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"t1", "1-5m"},
		{"t2", "5-m"},
		{"t3", "6-h"},
		{"t4", "7-D"},
		{"t5", "8-s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := GetOrParseNew(tt.name, tt.expr)
			t.Logf("String() raw = %v | %v, got = %v", tt.name, tt.expr, rl)
		})
	}
}
