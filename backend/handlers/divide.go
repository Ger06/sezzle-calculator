package handlers

import (
	"net/http"

	"sezzle-calculator/backend/calculator"
)

type DivideRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

func Divide(w http.ResponseWriter, r *http.Request) {
	req, err := decodeRequest[DivideRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := requireField("a", req.A); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := requireField("b", req.B); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	result, err := calculator.Divide(*req.A, *req.B)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeResult(w, result)
}
