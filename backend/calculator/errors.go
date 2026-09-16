package calculator

import "errors"

var (
	ErrDivisionByZero  = errors.New("cannot divide by zero")
	ErrNegativeSqrt    = errors.New("cannot take the square root of a negative number")
	ErrNonFiniteResult = errors.New("result is not a finite real number")
)
