package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func get(t *testing.T, path string) (*http.Response, map[string]string) {
	t.Helper()
	srv := httptest.NewServer(NewMux())
	t.Cleanup(srv.Close)
	res, err := http.Get(srv.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { res.Body.Close() })
	body := map[string]string{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return res, body
}

func TestHealth(t *testing.T) {
	res, body := get(t, "/healthz")
	if res.StatusCode != http.StatusOK || body["status"] != "ok" {
		t.Errorf("healthz = %d %v", res.StatusCode, body)
	}
}

func TestHello(t *testing.T) {
	res, body := get(t, "/api/hello/Go")
	if res.StatusCode != http.StatusOK || body["message"] != "Hello, Go!" {
		t.Errorf("hello = %d %v", res.StatusCode, body)
	}
}
