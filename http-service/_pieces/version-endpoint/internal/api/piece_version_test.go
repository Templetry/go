package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersionEndpoint(t *testing.T) {
	srv := httptest.NewServer(NewMux())
	defer srv.Close()
	res, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body := map[string]string{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK || body["version"] != version {
		t.Errorf("version = %d %v, want %q", res.StatusCode, body, version)
	}
}
