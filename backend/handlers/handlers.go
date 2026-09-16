package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"sezzle-calculator/backend/calculator"
)

func decodeRequest[T any](r *http.Request) (T, error) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func requireField(name string, p *float64) error {
	if p == nil {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

func writeResult(w http.ResponseWriter, result float64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Result float64 `json:"result"`
	}{Result: result})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{Error: struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{Code: code, Message: message}})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", err.Error())
	case errors.Is(err, calculator.ErrNegativeSqrt):
		writeError(w, http.StatusUnprocessableEntity, "NEGATIVE_SQRT", err.Error())
	case errors.Is(err, calculator.ErrNonFiniteResult):
		writeError(w, http.StatusUnprocessableEntity, "NON_FINITE_RESULT", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected error")
	}
}
