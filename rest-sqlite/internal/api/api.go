// Package api wires the HTTP routes over the store.
package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"example.com/template-app/internal/store"
)

// Server holds the dependencies handlers need.
type Server struct {
	DB    *sql.DB
	Notes store.Notes
}

// registrars is the piece socket: pieces add their routes from their own
// file's init, so NewMux never needs editing.
var registrars []func(*http.ServeMux, *Server)

// Register adds a route registrar. Call it from an init function.
func Register(f func(*http.ServeMux, *Server)) { registrars = append(registrars, f) }

// options is what NewMux was configured with.
type options struct {
	environment string
}

// Option configures the mux. Variadic on purpose: NewMux(db) keeps
// compiling, so the piece socket and existing callers need no change.
type Option func(*options)

// tpl:if environments

// WithEnvironment makes /healthz report the active profile.
func WithEnvironment(name string) Option {
	return func(o *options) { o.environment = name }
}

// tpl:endif

// NewMux builds the router for a database handle.
func NewMux(db *sql.DB, opts ...Option) *http.ServeMux {
	var cfg options
	for _, apply := range opts {
		apply(&cfg)
	}

	s := &Server{DB: db, Notes: store.Notes{DB: db}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		body := map[string]string{"status": "ok"}
		if cfg.environment != "" {
			body["environment"] = cfg.environment
		}
		WriteJSON(w, http.StatusOK, body)
	})
	mux.HandleFunc("POST /api/notes", s.createNote)
	mux.HandleFunc("GET /api/notes", s.listNotes)
	mux.HandleFunc("GET /api/notes/{id}", s.getNote)
	mux.HandleFunc("PUT /api/notes/{id}", s.updateNote)
	mux.HandleFunc("DELETE /api/notes/{id}", s.deleteNote)
	for _, register := range registrars {
		register(mux, s)
	}
	return mux
}

// WriteJSON writes a JSON response; pieces use it too.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// Fail maps store errors to status codes.
func Fail(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

// PathID reads the {id} path value.
func PathID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

type noteInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	var in noteInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Title == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}
	n, err := s.Notes.Create(in.Title, in.Body)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, n)
}

func (s *Server) listNotes(w http.ResponseWriter, r *http.Request) {
	limit, offset := 50, 0
	if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v > 0 {
		offset = v
	}
	notes, err := s.Notes.List(limit, offset)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, notes)
}

func (s *Server) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := PathID(r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	n, err := s.Notes.Get(id)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, n)
}

func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {
	id, err := PathID(r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var in noteInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Title == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}
	n, err := s.Notes.Update(id, in.Title, in.Body)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, n)
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := PathID(r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := s.Notes.Delete(id); err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusNoContent, nil)
}
