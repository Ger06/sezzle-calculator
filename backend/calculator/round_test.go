package calculator

import "testing"

func TestRound(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"integer value", 5, 5},
		{"simple decimal", 2.5, 2.5},
		{"floating-point noise from arithmetic", 0.1 + 0.2, 0.3},
		{"more than 10 decimal places rounds down", 1.00000000004, 1},
		{"more than 10 decimal places rounds up", 1.00000000006, 1.0000000001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Round(tt.input)
			if got != tt.want {
				t.Errorf("Round(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
