package ngamux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestGetQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
	result := GetQuery(req, "id")

	if result != "1" {
		t.Errorf("TestGetQuery need %v, but got %v", "1", result)
	}

	result = GetQuery(req, "slug", "undefined")
	if result != "undefined" {
		t.Errorf("TestGetQuery need %v, but got %v", "undefined", result)
	}

	result = GetQuery(req, "slug")
	if result != "" {
		t.Errorf("TestGetQuery need %v, but got %v", "\"\"", result)
	}
}

func TestGetJSON(t *testing.T) {
	input := strings.NewReader(`{"id": 1}`)
	req := httptest.NewRequest(http.MethodGet, "/", input)

	var data map[string]interface{}
	err := GetJSON(req, &data)
	if err != nil {
		t.Errorf("TestGetJSON need %v, but got %v", "nil", err)
	}

	if data["id"] == nil {
		t.Errorf("TestGetJSON need %v, but got %v", "value", data["id"])
	}

	id, ok := data["id"].(float64)
	if !ok {
		t.Errorf("TestGetJSON need %v, but got %v", "true", ok)
	}

	if id != 1 {
		t.Errorf("TestGetJSON need %v, but got %v", 1, id)
	}
}

func TestSetContextValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = SetContextValue(req, "id", 1)

	id := req.Context().Value("id")
	if id != 1 {
		t.Errorf("TestGetJSON need %v, but got %v", 1, id)
	}

	slug := req.Context().Value("slug")
	if id != 1 {
		t.Errorf("TestGetJSON need %v, but got %v", nil, slug)
	}
}
