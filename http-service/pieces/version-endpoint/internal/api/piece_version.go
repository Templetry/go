package api

import "net/http"

// version is reported by GET /version. Release builds can override it with
// -ldflags "-X example.com/template-app/internal/api.version=...".
var version = "0.1.0" // tpl:var version_value 0.1.0

// init plugs this piece into the service through the api socket — no
// existing file is modified (ADR-0014).
func init() {
	Register(func(mux *http.ServeMux) {
		mux.HandleFunc("GET /version", func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"version": version})
		})
	})
}
