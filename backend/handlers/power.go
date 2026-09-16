package handlers

import (
	"net/http"

	"sezzle-calculator/backend/calculator"
)

type PowerRequest struct {
	Base     *float64 `json:"base"`
	Exponent *float64 `json:"exponent"`
}

func Power(w http.ResponseWriter, r *http.Request) {
	req, err := decodeRequest[PowerRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := requireField("base", req.Base); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := requireField("exponent", req.Exponent); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	result, err := calculator.Power(*req.Base, *req.Exponent)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeResult(w, result)
}
