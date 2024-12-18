package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestBuildRoute(t *testing.T) {
	must := must.New(t)
	result := buildRoute("/", http.MethodGet, func(rw http.ResponseWriter, r *http.Request) {})
	expected := Route{
		Path:   "/",
		Method: http.MethodGet,
	}
	must.Equal(expected.Path, result.Path)
	must.Equal(expected.Method, result.Method)
}

func TestAddRoute(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.addRoute(buildRoute("/", http.MethodGet, nil))
	mux.addRoute(buildRoute("/a", http.MethodGet, nil))
	mux.addRoute(buildRoute("/b", http.MethodGet, nil))
	result := len(mux.routes)
	expected := 3
	must.Equal(expected, result)

	mux = New(WithLogLevel(LogLevelQuiet))
	mux.addRoute(buildRoute("/a/:a", http.MethodGet, nil))
	mux.addRoute(buildRoute("/b/:b", http.MethodGet, nil))
	mux.addRoute(buildRoute("/c/:c", http.MethodGet, nil))
	result = len(mux.routesParam)
	expected = 3
	must.Equal(expected, result)
}

func TestGetRoute(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler, req := mux.getRoute(req)
	handler.Handler(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)

	mux1 := New(WithLogLevel(LogLevelQuiet))
	mux1.Get("/:a", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	req1 := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec1 := httptest.NewRecorder()
	handler1, req1 := mux1.getRoute(req1)
	handler1.Handler(rec1, req1)

	result = strings.ReplaceAll(rec1.Body.String(), "\n", "")
	expected = "ok"
	must.Equal(expected, result)

	req2 := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec2 := httptest.NewRecorder()
	handler2, req2 := mux.getRoute(req2)
	handler2.Handler(rec2, req2)

	result = strings.ReplaceAll(rec2.Body.String(), "\n", "")
	expected = "not found"
	must.Equal(expected, result)

	req3 := httptest.NewRequest(http.MethodPost, "/", nil)
	rec3 := httptest.NewRecorder()
	handler2, req3 = mux.getRoute(req3)
	handler2.Handler(rec3, req3)

	result = strings.ReplaceAll(rec3.Body.String(), "\n", "")
	expected = "method not allowed"
	must.Equal(expected, result)
}

func BenchmarkRouter(b *testing.B) {
	mux := New(WithLogLevel(LogLevelQuiet))
	defaultHandler := func(rw http.ResponseWriter, r *http.Request) {
	}
	b.Run("add route", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mux.addRoute(buildRoute("/", http.MethodGet, defaultHandler))
		}
	})

	b.Run("get route", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mux.getRoute(httptest.NewRequest(http.MethodGet, "/", nil))
		}
	})
}
