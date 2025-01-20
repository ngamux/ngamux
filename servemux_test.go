package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestServeMuxUse(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	middleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
		}
	}
	mux.Use(middleware)
	mux.Use(middleware)
	mux.Use(middleware)

	result := len(mux.middlewares)
	expected := 3

	must.Equal(expected, result)
}

func TestServeMuxGet(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		must := must.New(t)
		mux := NewHttpServeMux()
		mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		must := must.New(t)
		mux := NewHttpServeMux()
		mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/cats", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "404 page not found"
		must.Equal(expected, result)
	})
}

// func TestHead(t *testing.T) {
// 	must := must.New(t)
// 	mux := New(
// 		WithLogLevel(LogLevelQuiet),
// 	)
// 	mux.Head("/", func(rw http.ResponseWriter, r *http.Request) {
// 		Res(rw).Text("ok")
// 	})
//
// 	rec := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodHead, "/", nil)
// 	mux.ServeHTTP(rec, req)
//
// 	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
// 	expected := ""
// 	must.Equal(expected, result)
// }

func TestServeMuxPost(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxPut(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Put("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxPatch(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Patch("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxDelete(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Delete("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxAll(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.All("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}
}
