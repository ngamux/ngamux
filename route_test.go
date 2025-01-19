package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestGetRoute(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler, pattern := mux.Handler(req)
	must.NotNil(handler)
	must.Equal("/", pattern)
	handler.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)

	mux1 := New(WithLogLevel(LogLevelQuiet))
	mux1.Get("/{a}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	req1 := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec1 := httptest.NewRecorder()
	handler1, pattern1 := mux1.Handler(req1)
	must.NotNil(handler1)
	must.Equal("/{a}", pattern1)
	handler1.ServeHTTP(rec1, req1)

	result = strings.ReplaceAll(rec1.Body.String(), "\n", "")
	expected = "ok"
	must.Equal(expected, result)

	req2 := httptest.NewRequest(http.MethodGet, "/123", nil)
	handler2, pattern2 := mux.Handler(req2)
	must.Nil(handler2)
	must.Equal("", pattern2)
}
