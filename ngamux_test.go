package ngamux

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewNgamux(t *testing.T) {
	result := NewNgamux()
	expected := &Ngamux{
		routes:             routeMap{},
		routesParam:        routeMap{},
		config:             buildConfig(),
		regexpParamFounded: paramsFinder,
	}

	if !reflect.DeepEqual(result.ngamux().routes, expected.routes) {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.routes, result.ngamux().routes)
	}

	if !reflect.DeepEqual(result.ngamux().routesParam, expected.routesParam) {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.routesParam, result.ngamux().routesParam)
	}

	if result.ngamux().config.RemoveTrailingSlash != expected.config.RemoveTrailingSlash {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.config.RemoveTrailingSlash, result.ngamux().config.RemoveTrailingSlash)
	}

	if result.ngamux().regexpParamFounded != expected.regexpParamFounded {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.regexpParamFounded, result.ngamux().regexpParamFounded)
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

	result := len(mux.ngamux().middlewares)
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

func TestPatch(t *testing.T) {
	mux := NewNgamux()
	mux.Patch("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestPatch need %v, but got %v", expected, result)
	}
}

func TestDelete(t *testing.T) {
	mux := NewNgamux()
	mux.Delete("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestDelete need %v, but got %v", expected, result)
	}
}

func TestAll(t *testing.T) {
	mux := NewNgamux()
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

		if result != expected {
			t.Errorf("TestAll need %v, but got %v", expected, result)
		}
	}
}

func TestErrorResponse(t *testing.T) {
	mux := NewNgamux()
	mux.Get("/error-method", func(rw http.ResponseWriter, r *http.Request) error {
		return errors.New("something bad")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/error-method", nil)

	mux.ServeHTTP(rec, req)

	result := rec.Result()
	resBody := strings.ReplaceAll(rec.Body.String(), "\n", "")
	if resBody != "something bad" {
		t.Errorf("Expect body to \"something bad\", but got %s", resBody)
	}
	if result.StatusCode != 500 {
		t.Errorf("Status should be 500, but got %d", result.StatusCode)
	}
}
