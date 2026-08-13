// Package api wires the service routes (Go 1.22 method+path patterns).
package api

import (
	"encoding/json"
	"net/http"
)

// NewMux builds the service's router.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /api/hello/{name}", handleHello)
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	writeJSON(w, http.StatusOK, map[string]string{"message": "Hello, " + name + "!"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
