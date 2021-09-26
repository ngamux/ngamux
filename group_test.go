package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGroup(t *testing.T) {
	mux := NewNgamux()
	a := mux.Group("/a")
	a.Get("", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestGroup need %v, but got %v", expected, result)
	}
}
