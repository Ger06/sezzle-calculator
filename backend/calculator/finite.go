package calculator

import "math"

func checkFinite(result float64) error {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return ErrNonFiniteResult
	}
	return nil
}
