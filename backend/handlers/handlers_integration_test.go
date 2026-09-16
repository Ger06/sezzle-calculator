package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /add", Add)
	mux.HandleFunc("POST /subtract", Subtract)
	mux.HandleFunc("POST /multiply", Multiply)
	mux.HandleFunc("POST /divide", Divide)
	mux.HandleFunc("POST /sqrt", Sqrt)
	mux.HandleFunc("POST /power", Power)
	mux.HandleFunc("POST /percentage", Percentage)
	return CORSMiddleware(mux)
}

func TestIntegrationAllSevenOperationsSucceed(t *testing.T) {
	router := newRouter()

	tests := []struct {
		name   string
		path   string
		body   string
		result float64
	}{
		{"add", "/add", `{"a": 2, "b": 3}`, 5},
		{"subtract", "/subtract", `{"a": 5, "b": 3}`, 2},
		{"multiply", "/multiply", `{"a": 4, "b": 2.5}`, 10},
		{"divide", "/divide", `{"a": 10, "b": 4}`, 2.5},
		{"sqrt", "/sqrt", `{"value": 16}`, 4},
		{"power", "/power", `{"base": 2, "exponent": 10}`, 1024},
		{"percentage", "/percentage", `{"value": 200, "percentage": 15}`, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Origin", "http://localhost:5173")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
			}
			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
			}
			var body struct {
				Result float64 `json:"result"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}
			if body.Result != tt.result {
				t.Errorf("result = %v, want %v", body.Result, tt.result)
			}
		})
	}
}

func TestIntegrationEveryDocumentedErrorCode(t *testing.T) {
	router := newRouter()

	tests := []struct {
		name       string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{"invalid input: missing field", "/add", `{"a": 2}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"invalid input: wrong JSON type", "/add", `{"a": "x", "b": 3}`, http.StatusBadRequest, "INVALID_INPUT"},
		{"invalid input: malformed JSON", "/add", `{"a": 2,`, http.StatusBadRequest, "INVALID_INPUT"},
		{"division by zero", "/divide", `{"a": 1, "b": 0}`, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO"},
		{"negative sqrt", "/sqrt", `{"value": -4}`, http.StatusUnprocessableEntity, "NEGATIVE_SQRT"},
		{"non-finite result", "/power", `{"base": -8, "exponent": 0.5}`, http.StatusUnprocessableEntity, "NON_FINITE_RESULT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			var body struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}
			if body.Error.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", body.Error.Code, tt.wantCode)
			}
			if body.Error.Message == "" {
				t.Error("error message should not be empty")
			}
		})
	}
}

func TestIntegrationOptionsPreflightForAllEndpoints(t *testing.T) {
	router := newRouter()
	paths := []string{"/add", "/subtract", "/multiply", "/divide", "/sqrt", "/power", "/percentage"}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodOptions, path, nil)
			req.Header.Set("Origin", "http://localhost:5173")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusNoContent {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
			}
			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
			}
		})
	}
}
