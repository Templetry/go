// Package api wires the service routes (Go 1.22 method+path patterns).
package api

import (
	"encoding/json"
	"net/http"
)

// registrars is the piece socket: decoupled units (Templetry pieces) add
// themselves here from their own file's init, so NewMux never needs editing.
var registrars []func(*http.ServeMux)

// Register adds a route registrar. Call it from an init function.
func Register(f func(*http.ServeMux)) {
	registrars = append(registrars, f)
}

// options is what NewMux was configured with.
type options struct {
	environment string
}

// Option configures the mux. Variadic on purpose: NewMux() keeps compiling,
// so the piece socket and existing callers need no change.
type Option func(*options)

// tpl:if environments

// WithEnvironment makes /healthz report the active profile.
func WithEnvironment(name string) Option {
	return func(o *options) { o.environment = name }
}

// tpl:endif

// NewMux builds the service's router.
func NewMux(opts ...Option) *http.ServeMux {
	var cfg options
	for _, apply := range opts {
		apply(&cfg)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth(cfg))
	mux.HandleFunc("GET /api/hello/{name}", handleHello)
	for _, register := range registrars {
		register(mux)
	}
	return mux
}

func handleHealth(cfg options) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		body := map[string]string{"status": "ok"}
		if cfg.environment != "" {
			body["environment"] = cfg.environment
		}
		writeJSON(w, http.StatusOK, body)
	}
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
