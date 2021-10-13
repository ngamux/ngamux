package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildRoute(t *testing.T) {
	result := buildRoute("/", http.MethodGet, func(rw http.ResponseWriter, r *http.Request) error { return nil })
	expected := Route{
		Path:   "/",
		Method: http.MethodGet,
	}

	if result.Path != expected.Path {
		t.Errorf("TestBuildRoute need %v, but got %v", expected.Path, result.Path)
	}

	if result.Method != expected.Method {
		t.Errorf("TestBuildRoute need %v, but got %v", expected.Method, result.Method)
	}
}

func TestAddRoute(t *testing.T) {
	mux := NewNgamux()
	mux.addRoute(buildRoute("/", http.MethodGet, nil))
	mux.addRoute(buildRoute("/a", http.MethodGet, nil))
	mux.addRoute(buildRoute("/b", http.MethodGet, nil))
	result := len(mux.ngamux().routes)
	expected := 3

	if result != expected {
		t.Errorf("TestAddRoute need %v, but got %v", expected, result)
	}

	mux = NewNgamux()
	mux.addRoute(buildRoute("/a/:a", http.MethodGet, nil))
	mux.addRoute(buildRoute("/b/:b", http.MethodGet, nil))
	mux.addRoute(buildRoute("/c/:c", http.MethodGet, nil))
	result = len(mux.ngamux().routesParam)
	expected = 3

	if result != expected {
		t.Errorf("TestAddRoute need %v, but got %v", expected, result)
	}
}

func TestGetRoute(t *testing.T) {
	mux := NewNgamux()
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler, req := mux.getRoute(req)
	handler.Handler(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"

	if result != expected {
		t.Errorf("TestGetRoute need %v, but got %v", expected, result)
	}

	mux1 := NewNgamux()
	mux1.Get("/:a", func(rw http.ResponseWriter, r *http.Request) error {
		return String(rw, "ok")
	})

	req1 := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec1 := httptest.NewRecorder()
	handler1, req1 := mux1.getRoute(req1)
	handler1.Handler(rec1, req1)

	result = strings.ReplaceAll(rec1.Body.String(), "\n", "")
	expected = "ok"

	if result != expected {
		t.Errorf("TestGetRoute need %v, but got %v", expected, result)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec2 := httptest.NewRecorder()
	handler2, req2 := mux.getRoute(req2)
	handler2.Handler(rec2, req2)

	result = strings.ReplaceAll(rec2.Body.String(), "\n", "")
	expected = "not found"

	if result != expected {
		t.Errorf("TestGetRoute need %v, but got %v", expected, result)
	}

	req3 := httptest.NewRequest(http.MethodPost, "/", nil)
	rec3 := httptest.NewRecorder()
	handler2, req3 = mux.getRoute(req3)
	handler2.Handler(rec3, req3)

	result = strings.ReplaceAll(rec3.Body.String(), "\n", "")
	expected = "method not allowed"

	if result != expected {
		t.Errorf("TestGetRoute need %v, but got %v", expected, result)
	}
}
