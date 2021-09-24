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
