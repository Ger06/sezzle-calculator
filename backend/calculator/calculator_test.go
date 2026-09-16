package calculator

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"two positives", 2, 3, 5},
		{"negative first operand", -5, 3, -2},
		{"negative second operand", 5, -3, 2},
		{"both negative", -5, -3, -8},
		{"zero operands", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Add(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
			}
			if got != tt.want {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"two positives", 5, 3, 2},
		{"negative first operand", -5, 3, -8},
		{"negative second operand", 5, -3, 8},
		{"both negative", -5, -3, -2},
		{"zero operands", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subtract(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Subtract(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
			}
			if got != tt.want {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"two positives", 4, 2.5, 10},
		{"negative first operand", -4, 2, -8},
		{"negative second operand", 4, -2, -8},
		{"both negative", -4, -2, 8},
		{"zero operand", 0, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Multiply(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Multiply(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
			}
			if got != tt.want {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"two positives", 10, 4, 2.5, nil},
		{"negative dividend", -10, 4, -2.5, nil},
		{"negative divisor", 10, -4, -2.5, nil},
		{"both negative", -10, -4, 2.5, nil},
		{"divisor zero", 5, 0, 0, ErrDivisionByZero},
		{"dividend and divisor both zero", 0, 0, 0, ErrDivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("Divide(%v, %v) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Divide(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
			}
			if got != tt.want {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		want    float64
		wantErr error
	}{
		{"perfect square", 16, 4, nil},
		{"zero", 0, 0, nil},
		{"negative zero is not negative", math.Copysign(0, -1), 0, nil},
		{"non-perfect-square value", 2, Round(math.Sqrt(2)), nil},
		{"floating-point-boundary negative value", -1e-15, 0, ErrNegativeSqrt},
		{"clearly negative value", -4, 0, ErrNegativeSqrt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sqrt(tt.value)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("Sqrt(%v) error = %v, want %v", tt.value, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Sqrt(%v) returned unexpected error: %v", tt.value, err)
			}
			if got != tt.want {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		name           string
		base, exponent float64
		want           float64
		wantErr        error
	}{
		{"positive base, positive integer exponent", 2, 10, 1024, nil},
		{"any base, zero exponent", 5, 0, 1, nil},
		{"zero base, zero exponent", 0, 0, 1, nil},
		{"negative base, positive integer exponent", -2, 3, -8, nil},
		{"positive base, negative exponent", 2, -1, 0.5, nil},
		{"negative base, non-integer exponent", -8, 0.5, 0, ErrNonFiniteResult},
		{"zero base, negative exponent", 0, -1, 0, ErrNonFiniteResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Power(tt.base, tt.exponent)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("Power(%v, %v) error = %v, want %v", tt.base, tt.exponent, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Power(%v, %v) returned unexpected error: %v", tt.base, tt.exponent, err)
			}
			if got != tt.want {
				t.Errorf("Power(%v, %v) = %v, want %v", tt.base, tt.exponent, got, tt.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		name              string
		value, percentage float64
		want              float64
	}{
		{"typical percentage", 200, 15, 30},
		{"negative value", -50, 10, -5},
		{"negative percentage", 50, -10, -5},
		{"both negative", -50, -10, 5},
		{"zero percentage", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Percentage(tt.value, tt.percentage)
			if err != nil {
				t.Fatalf("Percentage(%v, %v) returned unexpected error: %v", tt.value, tt.percentage, err)
			}
			if got != tt.want {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tt.value, tt.percentage, got, tt.want)
			}
		})
	}
}
