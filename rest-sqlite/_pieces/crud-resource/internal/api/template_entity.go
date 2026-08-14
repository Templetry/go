package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/template-app/internal/store"
)

// init plugs this resource's routes into the mux through the api socket.
func init() {
	Register(func(mux *http.ServeMux, s *Server) {
		repo := store.TemplateEntities{DB: s.DB}
		h := templateEntityHandlers{repo: repo}
		mux.HandleFunc("POST /api/template-entities", h.create)
		mux.HandleFunc("GET /api/template-entities", h.list)
		mux.HandleFunc("GET /api/template-entities/{id}", h.get)
		mux.HandleFunc("PUT /api/template-entities/{id}", h.update)
		mux.HandleFunc("DELETE /api/template-entities/{id}", h.remove)
	})
}

type templateEntityHandlers struct {
	repo store.TemplateEntities
}

type templateEntityInput struct {
	Name  string `json:"name"`
	Notes string `json:"notes"`
}

func (h templateEntityHandlers) create(w http.ResponseWriter, r *http.Request) {
	var in templateEntityInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}
	row, err := h.repo.Create(in.Name, in.Notes)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, row)
}

func (h templateEntityHandlers) list(w http.ResponseWriter, r *http.Request) {
	limit, offset := 50, 0
	if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v > 0 {
		offset = v
	}
	rows, err := h.repo.List(limit, offset)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, rows)
}

func (h templateEntityHandlers) get(w http.ResponseWriter, r *http.Request) {
	id, err := PathID(r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	row, err := h.repo.Get(id)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, row)
}

func (h templateEntityHandlers) update(w http.ResponseWriter, r *http.Request) {
	id, err := PathID(r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var in templateEntityInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}
	row, err := h.repo.Update(id, in.Name, in.Notes)
	if err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, row)
}

func (h templateEntityHandlers) remove(w http.ResponseWriter, r *http.Request) {
	id, err := PathID(r)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.repo.Delete(id); err != nil {
		Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusNoContent, nil)
}
