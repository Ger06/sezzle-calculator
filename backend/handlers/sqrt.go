package handlers

import (
	"net/http"

	"sezzle-calculator/backend/calculator"
)

type SqrtRequest struct {
	Value *float64 `json:"value"`
}

func Sqrt(w http.ResponseWriter, r *http.Request) {
	req, err := decodeRequest[SqrtRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := requireField("value", req.Value); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	result, err := calculator.Sqrt(*req.Value)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeResult(w, result)
}
