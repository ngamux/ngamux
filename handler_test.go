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
	req := httptest.NewRequest(http.MethodGet, "/a", nil)

	{
		rec := httptest.NewRecorder()
		handler := ToHandler(func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		})
		handler.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}

	{
		rec := httptest.NewRecorder()
		handler := ToHandlerFunc(ToHandler(func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		}))
		handler.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}
}
