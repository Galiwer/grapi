package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	api := New("1.1.0", "test")

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

	if result["status"] != "healthy" {
		t.Errorf("expected healthy, got %v", result["status"])
	}
}

func TestVersion(t *testing.T) {
	api := New("1.1.0", "test")

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

	if result["version"] != "1.1.0" {
		t.Errorf("expected version 1.1.0, got %v", result["version"])
	}
}

func TestMetrics(t *testing.T) {
	api := New("1.1.0", "test")

	api.Home(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	api.Metrics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `grapi_http_requests_total{endpoint="/"} 1`) {
		t.Errorf("expected home counter in metrics output, got:\n%s", body)
	}
	if !strings.Contains(body, `grapi_info{version="1.1.0",environment="test"} 1`) {
		t.Errorf("expected grapi_info metric in output, got:\n%s", body)
	}
}
