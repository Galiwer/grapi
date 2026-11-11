package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

type API struct {
	AppVersion  string
	Environment string
	startTime   time.Time
	homeHits    atomic.Uint64
	healthHits  atomic.Uint64
	readyHits   atomic.Uint64
	versionHits atomic.Uint64
	messageHits atomic.Uint64
	metricsHits atomic.Uint64
}

func New(version, environment string) *API {
	return &API{
		AppVersion:  version,
		Environment: environment,
		startTime:   time.Now(),
	}
}

type Response map[string]interface{}

func writeJSON(w http.ResponseWriter, status int, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (a *API) Home(w http.ResponseWriter, r *http.Request) {
	a.homeHits.Add(1)
	writeJSON(w, http.StatusOK, Response{
		"application": "DevOps Demo API",
		"description": "REST API for demonstrating CI/CD and GitOps",
		"version":     a.AppVersion,
	})
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	a.healthHits.Add(1)
	writeJSON(w, http.StatusOK, Response{
		"status": "healthy",
	})
}

func (a *API) Ready(w http.ResponseWriter, r *http.Request) {
	a.readyHits.Add(1)
	writeJSON(w, http.StatusOK, Response{
		"status": "ready",
	})
}

func (a *API) Version(w http.ResponseWriter, r *http.Request) {
	a.versionHits.Add(1)
	writeJSON(w, http.StatusOK, Response{
		"application": "DevOps Demo API",
		"version":     a.AppVersion,
		"environment": a.Environment,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (a *API) Message(w http.ResponseWriter, r *http.Request) {
	a.messageHits.Add(1)
	writeJSON(w, http.StatusOK, Response{
		"message": "Hello! Your Go REST API is running successfully via automated GitOps CI/CD.",
		"status":  "success",
	})
}

func (a *API) Metrics(w http.ResponseWriter, r *http.Request) {
	a.metricsHits.Add(1)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	uptime := time.Since(a.startTime).Seconds()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "# HELP grapi_info Application version and environment metadata.\n")
	fmt.Fprintf(w, "# TYPE grapi_info gauge\n")
	fmt.Fprintf(w, "grapi_info{version=%q,environment=%q} 1\n", a.AppVersion, a.Environment)

	fmt.Fprintf(w, "# HELP grapi_uptime_seconds Time in seconds since the API server started.\n")
	fmt.Fprintf(w, "# TYPE grapi_uptime_seconds gauge\n")
	fmt.Fprintf(w, "grapi_uptime_seconds %.2f\n", uptime)

	fmt.Fprintf(w, "# HELP grapi_goroutines Number of active Go goroutines.\n")
	fmt.Fprintf(w, "# TYPE grapi_goroutines gauge\n")
	fmt.Fprintf(w, "grapi_goroutines %d\n", runtime.NumGoroutine())

	fmt.Fprintf(w, "# HELP grapi_memory_alloc_bytes Bytes of allocated heap objects.\n")
	fmt.Fprintf(w, "# TYPE grapi_memory_alloc_bytes gauge\n")
	fmt.Fprintf(w, "grapi_memory_alloc_bytes %d\n", mem.Alloc)

	fmt.Fprintf(w, "# HELP grapi_http_requests_total Total number of HTTP requests processed by endpoint.\n")
	fmt.Fprintf(w, "# TYPE grapi_http_requests_total counter\n")
	fmt.Fprintf(w, "grapi_http_requests_total{endpoint=\"/\"} %d\n", a.homeHits.Load())
	fmt.Fprintf(w, "grapi_http_requests_total{endpoint=\"/health\"} %d\n", a.healthHits.Load())
	fmt.Fprintf(w, "grapi_http_requests_total{endpoint=\"/ready\"} %d\n", a.readyHits.Load())
	fmt.Fprintf(w, "grapi_http_requests_total{endpoint=\"/version\"} %d\n", a.versionHits.Load())
	fmt.Fprintf(w, "grapi_http_requests_total{endpoint=\"/api/message\"} %d\n", a.messageHits.Load())
	fmt.Fprintf(w, "grapi_http_requests_total{endpoint=\"/metrics\"} %d\n", a.metricsHits.Load())
}
