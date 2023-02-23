package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestServeHTTP(t *testing.T) {
	must := must.New(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a", nil)
	handler := Handler(func(rw http.ResponseWriter, r *http.Request) error {
		return Res(rw).String("ok")
	})
	handler.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}
