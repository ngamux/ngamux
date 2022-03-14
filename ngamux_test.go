package ngamux

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestNewNgamux(t *testing.T) {
	must := must.New(t)
	result := New()
	expected := &Ngamux{
		routes:            routeMap{},
		routesParam:       routeMap{},
		config:            buildConfig(),
		regexpParamFinded: paramsFinder,
	}

	must.Equal(expected.routes, result.routes)
	must.Equal(expected.routesParam, result.routesParam)
	must.Equal(expected.config.RemoveTrailingSlash, result.config.RemoveTrailingSlash)
	must.Equal(expected.regexpParamFinded, result.regexpParamFinded)
}

func TestUse(t *testing.T) {
	must := must.New(t)
	mux := New()
	middleware := func(next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) error {
			return nil
		}
	}
	mux.Use(middleware)
	mux.Use(middleware)
	mux.Use(middleware)

	result := len(mux.middlewares)
	expected := 3

	must.Equal(expected, result)
}

func TestGet(t *testing.T) {
	must := must.New(t)
	mux := New()
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestPost(t *testing.T) {
	must := must.New(t)
	mux := New()
	mux.Post("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestPut(t *testing.T) {
	must := must.New(t)
	mux := New()
	mux.Put("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestPatch(t *testing.T) {
	must := must.New(t)
	mux := New()
	mux.Patch("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestDelete(t *testing.T) {
	must := must.New(t)
	mux := New()
	mux.Delete("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestAll(t *testing.T) {
	must := must.New(t)
	mux := New()
	mux.All("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
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

func TestErrorResponse(t *testing.T) {
	must := must.New(t)
	mux := New()
	mux.Get("/error-method", func(rw http.ResponseWriter, r *http.Request) error {
		return errors.New("something bad")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/error-method", nil)

	mux.ServeHTTP(rec, req)

	result := rec.Result()
	resultBody := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expectedBody := "something bad"
	must.Equal(expectedBody, resultBody)

	resultStatus := result.StatusCode
	expectedStatus := 500
	must.Equal(expectedStatus, resultStatus)
}
