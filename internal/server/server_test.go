package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postJSON(t *testing.T, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	return rec
}

func TestEqEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/eq", map[string]interface{}{
		"alpha": 1.1, "beta": 0.4, "gamma": 0.4, "delta": 0.1,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestOrbitEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/orbit", map[string]interface{}{
		"alpha": 1.1, "beta": 0.4, "gamma": 0.4, "delta": 0.1,
		"V0": 1.2, "P0": 2, "t_end": 10, "steps": 500,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestInvalidEqReturns400(t *testing.T) {
	rec := postJSON(t, "/api/eq", map[string]interface{}{
		"alpha": 0, "beta": 0.4, "gamma": 0.4, "delta": 0.1,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/eq", nil)
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d", rec.Code)
	}
}
