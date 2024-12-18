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
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	handler := func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	}
	a := mux.Group("/a")
	a.Get("", handler)
	a.Post("", handler)
	a.Put("", handler)
	a.Patch("", handler)
	a.Delete("", handler)
	a.All("", handler)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}
