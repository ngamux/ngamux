package ngamux

import (
	"net/http"
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
	result := len(mux.routes)
	expected := 3

	if result != expected {
		t.Errorf("TestAddRoute need %v, but got %v", expected, result)
	}

	mux = NewNgamux()
	mux.addRoute(buildRoute("/a/:a", http.MethodGet, nil))
	mux.addRoute(buildRoute("/b/:b", http.MethodGet, nil))
	mux.addRoute(buildRoute("/c/:c", http.MethodGet, nil))
	result = len(mux.routesParam)
	expected = 3

	if result != expected {
		t.Errorf("TestAddRoute need %v, but got %v", expected, result)
	}
}

