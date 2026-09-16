package calculator

import "math"

func Add(a, b float64) (float64, error) {
	result := a + b
	if err := checkFinite(result); err != nil {
		return 0, err
	}
	return Round(result), nil
}

func Subtract(a, b float64) (float64, error) {
	result := a - b
	if err := checkFinite(result); err != nil {
		return 0, err
	}
	return Round(result), nil
}

func Multiply(a, b float64) (float64, error) {
	result := a * b
	if err := checkFinite(result); err != nil {
		return 0, err
	}
	return Round(result), nil
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	result := a / b
	if err := checkFinite(result); err != nil {
		return 0, err
	}
	return Round(result), nil
}

func Sqrt(value float64) (float64, error) {
	if value < 0 {
		return 0, ErrNegativeSqrt
	}
	result := math.Sqrt(value)
	if err := checkFinite(result); err != nil {
		return 0, err
	}
	return Round(result), nil
}

func Power(base, exponent float64) (float64, error) {
	result := math.Pow(base, exponent)
	if err := checkFinite(result); err != nil {
		return 0, err
	}
	return Round(result), nil
}

func Percentage(value, percentage float64) (float64, error) {
	result := value * percentage / 100
	if err := checkFinite(result); err != nil {
		return 0, err
	}
	return Round(result), nil
}
