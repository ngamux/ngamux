package ngamux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestWithMiddlewares(t *testing.T) {
	must := must.New(t)
	result := WithMiddlewares()(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	must.NotNil(result)

	result = WithMiddlewares(nil)(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	must.NotNil(result)

	result = WithMiddlewares(nil)(nil)
	must.Nil(result)
}

func TestGetParam(t *testing.T) {
	must := must.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), KeyContextParams, [][]string{{"id", "1"}}))
	result := GetParam(req, "id")
	must.Equal("1", result)

	result = GetParam(req, "slug")
	must.Equal("", result)
}

func TestGetQuery(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
	result := GetQuery(req, "id")
	must.Equal("1", result)

	result = GetQuery(req, "slug", "undefined")
	must.Equal("undefined", result)

	result = GetQuery(req, "slug")
	must.Equal("", result)
}

func TestGetJSON(t *testing.T) {
	must := must.New(t)
	input := strings.NewReader(`{"id": 1}`)
	req := httptest.NewRequest(http.MethodGet, "/", input)

	var data map[string]any
	err := GetJSON(req, &data)
	must.Nil(err)
	must.NotNil(data["id"])

	id, ok := data["id"].(float64)
	must.True(ok)
	must.Equal(float64(1), id)
}

func TestGetFormValue(t *testing.T) {

	t.Run("can supply default value", func(t *testing.T) {
		must := must.New(t)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		result := GetFormValue(req, "idontknow", "1")
		must.Equal("1", result)
	})
}

func TestSetContextValue(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = SetContextValue(req, "id", 1)

	id := req.Context().Value("id")
	must.Equal(1, id)

	slug := req.Context().Value("slug")
	must.Nil(slug)
}

type Key string

const KeyID Key = "id"

func TestGetContextValue(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), KeyID, 1))

	id := GetContextValue(req, KeyID)
	must.Equal(1, id)

	slug := GetContextValue(req, "slug")
	must.Nil(slug)
}
