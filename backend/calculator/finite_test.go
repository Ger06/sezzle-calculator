package calculator

import (
	"math"
	"testing"
)

func TestCheckFinite(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantErr error
	}{
		{"finite positive value", 42.5, nil},
		{"zero", 0, nil},
		{"finite negative value", -17, nil},
		{"NaN", math.NaN(), ErrNonFiniteResult},
		{"positive infinity", math.Inf(1), ErrNonFiniteResult},
		{"negative infinity", math.Inf(-1), ErrNonFiniteResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkFinite(tt.input)
			if tt.wantErr == nil && err != nil {
				t.Errorf("checkFinite(%v) = %v, want nil", tt.input, err)
			}
			if tt.wantErr != nil && err != tt.wantErr {
				t.Errorf("checkFinite(%v) = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
