package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a", nil)
	handler := Handler(func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})
	handler.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestServeHTTP need %v, but got %v", expected, result)
	}
}
