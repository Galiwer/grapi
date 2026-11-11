package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type API struct {
	AppVersion  string
	Environment string
}

func New(version, environment string) *API {
	return &API{
		AppVersion:  version,
		Environment: environment,
	}
}

type Response map[string]interface{}

func writeJSON(w http.ResponseWriter, status int, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (a *API) Home(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Response{
		"application": "DevOps Demo API",
		"description": "REST API for demonstrating CI/CD and GitOps",
		"version":     a.AppVersion,
	})
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Response{
		"status": "healthy",
	})
}

func (a *API) Ready(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Response{
		"status": "ready",
	})
}

func (a *API) Version(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Response{
		"application": "DevOps Demo API",
		"version":     a.AppVersion,
		"environment": a.Environment,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (a *API) Message(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Response{
		"message": "Hello! Your Go REST API is running successfully.",
		"status":  "success",
	})
}
