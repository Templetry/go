package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/template-app/internal/store"
)

// newTestServer boots the API over a throwaway in-memory database.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Migrate(db); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewMux(db))
	t.Cleanup(func() { srv.Close(); db.Close() })
	return srv
}

func do(t *testing.T, method, url string, body any) (*http.Response, []byte) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	out := make([]byte, 0)
	buf.Reset()
	_, _ = buf.ReadFrom(res.Body)
	out = buf.Bytes()
	return res, out
}

func TestHealth(t *testing.T) {
	srv := newTestServer(t)
	res, body := do(t, "GET", srv.URL+"/healthz", nil)
	if res.StatusCode != http.StatusOK || !bytes.Contains(body, []byte("ok")) {
		t.Errorf("healthz = %d %s", res.StatusCode, body)
	}
}

func TestNotesCRUD(t *testing.T) {
	srv := newTestServer(t)

	res, body := do(t, "POST", srv.URL+"/api/notes", map[string]string{"title": "first", "body": "b"})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create = %d %s", res.StatusCode, body)
	}
	var created store.Note
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 || created.Title != "first" || created.CreatedAt == "" {
		t.Fatalf("created = %+v", created)
	}

	res, body = do(t, "GET", srv.URL+"/api/notes", nil)
	var list []store.Note
	if err := json.Unmarshal(body, &list); err != nil || len(list) != 1 {
		t.Fatalf("list = %d %s (%v)", res.StatusCode, body, err)
	}

	res, _ = do(t, "PUT", srv.URL+"/api/notes/1", map[string]string{"title": "second"})
	if res.StatusCode != http.StatusOK {
		t.Errorf("update = %d", res.StatusCode)
	}

	res, _ = do(t, "DELETE", srv.URL+"/api/notes/1", nil)
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("delete = %d", res.StatusCode)
	}
	res, _ = do(t, "GET", srv.URL+"/api/notes/1", nil)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("get after delete = %d, want 404", res.StatusCode)
	}
}

func TestValidation(t *testing.T) {
	srv := newTestServer(t)
	res, _ := do(t, "POST", srv.URL+"/api/notes", map[string]string{"body": "no title"})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("create without title = %d, want 400", res.StatusCode)
	}
	res, _ = do(t, "GET", srv.URL+"/api/notes/abc", nil)
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("bad id = %d, want 400", res.StatusCode)
	}
}
