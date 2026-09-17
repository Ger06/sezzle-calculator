package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

func TestDecodeRequestValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"a": 2, "b": 3}`))

	got, err := decodeRequest[testRequest](req)
	if err != nil {
		t.Fatalf("decodeRequest returned unexpected error: %v", err)
	}
	if got.A == nil || *got.A != 2 {
		t.Errorf("A = %v, want 2", got.A)
	}
	if got.B == nil || *got.B != 3 {
		t.Errorf("B = %v, want 3", got.B)
	}
}

func TestDecodeRequestMissingFieldDecodesAsNil(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"a": 2}`))

	got, err := decodeRequest[testRequest](req)
	if err != nil {
		t.Fatalf("decodeRequest returned unexpected error: %v", err)
	}
	if got.B != nil {
		t.Errorf("B = %v, want nil for an absent field", *got.B)
	}
}

func TestDecodeRequestWrongJSONType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"a": "not-a-number", "b": 3}`))

	_, err := decodeRequest[testRequest](req)
	if err == nil {
		t.Fatal("decodeRequest returned nil error for a wrong JSON type, want an error")
	}
	if !strings.Contains(err.Error(), "cannot unmarshal") {
		t.Errorf("error message %q does not look like a JSON type error", err.Error())
	}
}

func TestDecodeRequestMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"a": 2,`))

	_, err := decodeRequest[testRequest](req)
	if err == nil {
		t.Fatal("decodeRequest returned nil error for malformed JSON, want an error")
	}
}

func TestRequireFieldMissing(t *testing.T) {
	err := requireField("b", nil)
	if err == nil {
		t.Fatal("requireField returned nil for a nil pointer, want an error")
	}
	want := "b is required"
	if err.Error() != want {
		t.Errorf("requireField error = %q, want %q", err.Error(), want)
	}
}

func TestRequireFieldPresentEvenWhenExplicitZero(t *testing.T) {
	value := 0.0
	if err := requireField("b", &value); err != nil {
		t.Errorf("requireField returned an error for an explicit 0 value: %v", err)
	}
}

func TestWrongTypeAndMissingFieldProduceDistinctMessages(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"a": "not-a-number", "b": 3}`))
	_, decodeErr := decodeRequest[testRequest](req)
	if decodeErr == nil {
		t.Fatal("expected a decode error for a wrong JSON type")
	}

	missingErr := requireField("b", nil)
	if missingErr == nil {
		t.Fatal("expected an error for a missing required field")
	}

	if decodeErr.Error() == missingErr.Error() {
		t.Errorf("wrong-type and missing-field errors should be distinct, both were %q", decodeErr.Error())
	}
}

func TestWriteResult(t *testing.T) {
	rec := httptest.NewRecorder()
	writeResult(rec, 5)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var body struct {
		Result float64 `json:"result"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if body.Result != 5 {
		t.Errorf("result = %v, want 5", body.Result)
	}
}

func TestWriteErrorValidation(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, http.StatusBadRequest, "INVALID_INPUT", "b is required")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
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
	if body.Error.Code != "INVALID_INPUT" {
		t.Errorf("code = %q, want %q", body.Error.Code, "INVALID_INPUT")
	}
	if body.Error.Message != "b is required" {
		t.Errorf("message = %q, want %q", body.Error.Message, "b is required")
	}
}

func TestWriteErrorDomain(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", "cannot divide by zero")

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
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
	if body.Error.Code != "DIVISION_BY_ZERO" {
		t.Errorf("code = %q, want %q", body.Error.Code, "DIVISION_BY_ZERO")
	}
	if body.Error.Message != "cannot divide by zero" {
		t.Errorf("message = %q, want %q", body.Error.Message, "cannot divide by zero")
	}
}

func decodeResult(t *testing.T, rec *httptest.ResponseRecorder) float64 {
	t.Helper()
	var body struct {
		Result float64 `json:"result"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	return body.Result
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) (code, message string) {
	t.Helper()
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	return body.Error.Code, body.Error.Message
}

func TestAddHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{"a": 2, "b": 3}`))
		rec := httptest.NewRecorder()

		Add(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := decodeResult(t, rec); got != 5 {
			t.Errorf("result = %v, want 5", got)
		}
	})

	t.Run("missing operand", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{"a": 2}`))
		rec := httptest.NewRecorder()

		Add(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		code, message := decodeError(t, rec)
		if code != "INVALID_INPUT" {
			t.Errorf("code = %q, want %q", code, "INVALID_INPUT")
		}
		if message != "b is required" {
			t.Errorf("message = %q, want %q", message, "b is required")
		}
	})
}

func TestSubtractHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/subtract", strings.NewReader(`{"a": 5, "b": 3}`))
		rec := httptest.NewRecorder()

		Subtract(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := decodeResult(t, rec); got != 2 {
			t.Errorf("result = %v, want 2", got)
		}
	})

	t.Run("missing operand", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/subtract", strings.NewReader(`{"b": 3}`))
		rec := httptest.NewRecorder()

		Subtract(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		code, message := decodeError(t, rec)
		if code != "INVALID_INPUT" {
			t.Errorf("code = %q, want %q", code, "INVALID_INPUT")
		}
		if message != "a is required" {
			t.Errorf("message = %q, want %q", message, "a is required")
		}
	})
}

func TestMultiplyHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/multiply", strings.NewReader(`{"a": 4, "b": 2.5}`))
		rec := httptest.NewRecorder()

		Multiply(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := decodeResult(t, rec); got != 10 {
			t.Errorf("result = %v, want 10", got)
		}
	})

	t.Run("missing operand", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/multiply", strings.NewReader(`{"a": 4}`))
		rec := httptest.NewRecorder()

		Multiply(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		code, message := decodeError(t, rec)
		if code != "INVALID_INPUT" {
			t.Errorf("code = %q, want %q", code, "INVALID_INPUT")
		}
		if message != "b is required" {
			t.Errorf("message = %q, want %q", message, "b is required")
		}
	})
}

func TestDivideHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/divide", strings.NewReader(`{"a": 10, "b": 4}`))
		rec := httptest.NewRecorder()

		Divide(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := decodeResult(t, rec); got != 2.5 {
			t.Errorf("result = %v, want 2.5", got)
		}
	})

	t.Run("missing operand", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/divide", strings.NewReader(`{"a": 10}`))
		rec := httptest.NewRecorder()

		Divide(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		code, message := decodeError(t, rec)
		if code != "INVALID_INPUT" {
			t.Errorf("code = %q, want %q", code, "INVALID_INPUT")
		}
		if message != "b is required" {
			t.Errorf("message = %q, want %q", message, "b is required")
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/divide", strings.NewReader(`{"a": 1, "b": 0}`))
		rec := httptest.NewRecorder()

		Divide(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
		}
		code, _ := decodeError(t, rec)
		if code != "DIVISION_BY_ZERO" {
			t.Errorf("code = %q, want %q", code, "DIVISION_BY_ZERO")
		}
	})
}

func TestSqrtHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sqrt", strings.NewReader(`{"value": 16}`))
		rec := httptest.NewRecorder()

		Sqrt(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := decodeResult(t, rec); got != 4 {
			t.Errorf("result = %v, want 4", got)
		}
	})

	t.Run("missing operand", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sqrt", strings.NewReader(`{}`))
		rec := httptest.NewRecorder()

		Sqrt(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		code, message := decodeError(t, rec)
		if code != "INVALID_INPUT" {
			t.Errorf("code = %q, want %q", code, "INVALID_INPUT")
		}
		if message != "value is required" {
			t.Errorf("message = %q, want %q", message, "value is required")
		}
	})

	t.Run("negative value", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sqrt", strings.NewReader(`{"value": -4}`))
		rec := httptest.NewRecorder()

		Sqrt(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
		}
		code, _ := decodeError(t, rec)
		if code != "NEGATIVE_SQRT" {
			t.Errorf("code = %q, want %q", code, "NEGATIVE_SQRT")
		}
	})
}

func TestPowerHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/power", strings.NewReader(`{"base": 2, "exponent": 10}`))
		rec := httptest.NewRecorder()

		Power(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := decodeResult(t, rec); got != 1024 {
			t.Errorf("result = %v, want 1024", got)
		}
	})

	t.Run("missing operand", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/power", strings.NewReader(`{"base": 2}`))
		rec := httptest.NewRecorder()

		Power(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		code, message := decodeError(t, rec)
		if code != "INVALID_INPUT" {
			t.Errorf("code = %q, want %q", code, "INVALID_INPUT")
		}
		if message != "exponent is required" {
			t.Errorf("message = %q, want %q", message, "exponent is required")
		}
	})

	t.Run("non-finite result", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/power", strings.NewReader(`{"base": -8, "exponent": 0.5}`))
		rec := httptest.NewRecorder()

		Power(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
		}
		code, _ := decodeError(t, rec)
		if code != "NON_FINITE_RESULT" {
			t.Errorf("code = %q, want %q", code, "NON_FINITE_RESULT")
		}
	})
}

func TestPercentageHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/percentage", strings.NewReader(`{"value": 200, "percentage": 15}`))
		rec := httptest.NewRecorder()

		Percentage(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if got := decodeResult(t, rec); got != 30 {
			t.Errorf("result = %v, want 30", got)
		}
	})

	t.Run("missing operand", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/percentage", strings.NewReader(`{"value": 200}`))
		rec := httptest.NewRecorder()

		Percentage(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		code, message := decodeError(t, rec)
		if code != "INVALID_INPUT" {
			t.Errorf("code = %q, want %q", code, "INVALID_INPUT")
		}
		if message != "percentage is required" {
			t.Errorf("message = %q, want %q", message, "percentage is required")
		}
	})
}

func TestCORSMiddlewareSetsHeadersAndCallsNext(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	CORSMiddleware(next).ServeHTTP(rec, req)

	if !nextCalled {
		t.Error("CORSMiddleware did not call the wrapped handler for a POST request")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCORSMiddlewareOptionsPreflight(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodOptions, "/add", nil)
	rec := httptest.NewRecorder()

	CORSMiddleware(next).ServeHTTP(rec, req)

	if nextCalled {
		t.Error("CORSMiddleware should short-circuit an OPTIONS preflight request, but called the wrapped handler")
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("Access-Control-Allow-Methods header missing on preflight response")
	}
}

func TestAllowedOriginRespectsEnvOverride(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGIN", "https://frontend-production.up.railway.app")

	req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	want := "https://frontend-production.up.railway.app"
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != want {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, want)
	}
}

func TestAllowedOriginFallsBackToLocalhostWhenUnset(t *testing.T) {
	// An empty value is indistinguishable from "unset" to os.Getenv, and
	// t.Setenv restores whatever CORS_ALLOWED_ORIGIN was after the test.
	t.Setenv("CORS_ALLOWED_ORIGIN", "")

	req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	want := "http://localhost:5173"
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != want {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, want)
	}
}
