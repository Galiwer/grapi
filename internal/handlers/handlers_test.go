package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	api := New("1.0.0", "test")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	api.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var result map[string]interface{}

	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	if result["status"] != "INTENTIONAL_CI_FAILURE_TEST" {
		t.Errorf("expected INTENTIONAL_CI_FAILURE_TEST, got %v", result["status"])
	}
}

func TestVersion(t *testing.T) {
	api := New("1.0.0", "test")

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	api.Version(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var result map[string]interface{}

	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	if result["version"] != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %v", result["version"])
	}
}
