package handlers

import (
	"net/http"

	"sezzle-calculator/backend/calculator"
)

type PercentageRequest struct {
	Value      *float64 `json:"value"`
	Percentage *float64 `json:"percentage"`
}

func Percentage(w http.ResponseWriter, r *http.Request) {
	req, err := decodeRequest[PercentageRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := requireField("value", req.Value); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := requireField("percentage", req.Percentage); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	result, err := calculator.Percentage(*req.Value, *req.Percentage)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeResult(w, result)
}
