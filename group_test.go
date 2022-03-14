package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestGroup(t *testing.T) {
	must := must.New(t)
	mux := New()
	a := mux.Group("/a")
	a.Get("", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}
