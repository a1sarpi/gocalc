package constants

import (
	"math"
	"testing"
)

func TestGetConstant(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     float64
		wantOk   bool
	}{
		{
			name:     "pi constant",
			constant: "pi",
			want:     math.Pi,
			wantOk:   true,
		},
		{
			name:     "e constant",
			constant: "e",
			want:     math.E,
			wantOk:   true,
		},
		{
			name:     "unknown constant",
			constant: "unknown",
			want:     0,
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetConstant(tt.constant)
			if ok != tt.wantOk {
				t.Errorf("GetConstant() ok = %v, want %v", ok, tt.wantOk)
				return
			}
			if ok && got != tt.want {
				t.Errorf("GetConstant() = %v, want %v", got, tt.want)
			}
		})
	}
} 