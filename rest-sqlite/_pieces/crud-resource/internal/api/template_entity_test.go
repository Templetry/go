package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"example.com/template-app/internal/store"
)

func TestTemplateEntityCRUD(t *testing.T) {
	srv := newTestServer(t)

	res, body := do(t, "POST", srv.URL+"/api/template-entities", map[string]string{"name": "first", "notes": "n"})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create = %d %s", res.StatusCode, body)
	}
	var created store.TemplateEntity
	if err := json.Unmarshal(body, &created); err != nil || created.ID == 0 || created.Name != "first" {
		t.Fatalf("created = %+v (%v)", created, err)
	}

	_, body = do(t, "GET", srv.URL+"/api/template-entities", nil)
	var list []store.TemplateEntity
	if err := json.Unmarshal(body, &list); err != nil || len(list) != 1 {
		t.Fatalf("list = %s (%v)", body, err)
	}

	res, _ = do(t, "PUT", srv.URL+"/api/template-entities/1", map[string]string{"name": "second"})
	if res.StatusCode != http.StatusOK {
		t.Errorf("update = %d", res.StatusCode)
	}

	res, _ = do(t, "DELETE", srv.URL+"/api/template-entities/1", nil)
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("delete = %d", res.StatusCode)
	}
	res, _ = do(t, "GET", srv.URL+"/api/template-entities/1", nil)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("get after delete = %d, want 404", res.StatusCode)
	}
}
