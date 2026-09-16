package calculator

import "math"

func Round(value float64) float64 {
	return math.Round(value*1e10) / 1e10
}
