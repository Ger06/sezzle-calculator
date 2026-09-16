package calculator

import (
	"errors"
	"testing"
)

func TestDomainErrorsAreDistinct(t *testing.T) {
	if errors.Is(ErrDivisionByZero, ErrNegativeSqrt) {
		t.Error("ErrDivisionByZero should not equal ErrNegativeSqrt")
	}
	if errors.Is(ErrDivisionByZero, ErrNonFiniteResult) {
		t.Error("ErrDivisionByZero should not equal ErrNonFiniteResult")
	}
	if errors.Is(ErrNegativeSqrt, ErrNonFiniteResult) {
		t.Error("ErrNegativeSqrt should not equal ErrNonFiniteResult")
	}
}

func TestDomainErrorsAreSelfEqual(t *testing.T) {
	if !errors.Is(ErrDivisionByZero, ErrDivisionByZero) {
		t.Error("ErrDivisionByZero should equal itself via errors.Is")
	}
	if !errors.Is(ErrNegativeSqrt, ErrNegativeSqrt) {
		t.Error("ErrNegativeSqrt should equal itself via errors.Is")
	}
	if !errors.Is(ErrNonFiniteResult, ErrNonFiniteResult) {
		t.Error("ErrNonFiniteResult should equal itself via errors.Is")
	}
}
