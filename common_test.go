package ngamux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWithMiddlewares(t *testing.T) {
	result := WithMiddlewares()(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	if result == nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", reflect.TypeOf(result), nil)
	}

	result = WithMiddlewares(nil)(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	if result == nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", reflect.TypeOf(result), nil)
	}

	result = WithMiddlewares(nil)(nil)
	if result != nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", nil, reflect.TypeOf(result))
	}
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), KeyContextParams, [][]string{{"id", "1"}}))
	result := GetParam(req, "id")

	if result != "1" {
		t.Errorf("TestGetParam need %v, but got %v", "1", result)
	}

	result = GetParam(req, "slug")
	if result != "" {
		t.Errorf("TestGetParam need %v, but got %v", "\"\"", result)
	}
}
