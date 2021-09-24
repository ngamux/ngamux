package ngamux

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewNgamux(t *testing.T) {
	result := NewNgamux()
	expected := &Ngamux{
		routes:            routeMap{},
		routesParam:       routeMap{},
		config:            buildConfig(),
		regexpParamFinded: paramsFinder,
	}

	if !reflect.DeepEqual(result.routes, expected.routes) {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.routes, result.routes)
	}

	if !reflect.DeepEqual(result.routesParam, expected.routesParam) {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.routesParam, result.routesParam)
	}

	if result.config.RemoveTrailingSlash != expected.config.RemoveTrailingSlash {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.config.RemoveTrailingSlash, result.config.RemoveTrailingSlash)
	}

	if result.regexpParamFinded != expected.regexpParamFinded {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.regexpParamFinded, result.regexpParamFinded)
	}
}

func TestUse(t *testing.T) {
	mux := NewNgamux()
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

	if result != expected {
		t.Errorf("TestUse need %v, but got %v", expected, result)
	}
}

func TestGet(t *testing.T) {
	mux := NewNgamux()
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestGet need %v, but got %v", expected, result)
	}
}

func TestPost(t *testing.T) {
	mux := NewNgamux()
	mux.Post("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestPost need %v, but got %v", expected, result)
	}
}

func TestPut(t *testing.T) {
	mux := NewNgamux()
	mux.Put("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestPut need %v, but got %v", expected, result)
	}
}

