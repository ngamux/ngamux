package ngamux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestReq(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	expected := &Request{req}

	result := Req(req)

	must.NotNil(result)
	must.Equal(expected, result)
}

func TestParams(t *testing.T) {
	must := must.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), KeyContextParams, [][]string{{"id", "1"}}))
	result := Req(req).Params("id")
	must.Equal("1", result)

	result = Req(req).Params("slug")
	must.Equal("", result)
}

func TestQuery(t *testing.T) {
	must := must.New(t)
	req := Req(httptest.NewRequest(http.MethodGet, "/?id=1", nil))
	result := req.Query("id")
	must.Equal("1", result)

	result = req.Query("slug", "undefined")
	must.Equal("undefined", result)

	result = req.Query("slug")
	must.Equal("", result)
}

func TestJSON(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		must := must.New(t)
		input := strings.NewReader(`{"id": 1}`)
		req := Req(httptest.NewRequest(http.MethodGet, "/?id=1", input))

		var data map[string]any
		err := req.JSON(&data)
		must.Nil(err)
		must.NotNil(data["id"])

		id, ok := data["id"].(float64)
		must.True(ok)
		must.Equal(float64(1), id)
	})

	t.Run("negative", func(t *testing.T) {
		must := must.New(t)
		input := strings.NewReader(`id`)
		req := Req(httptest.NewRequest(http.MethodGet, "/?id=1", input))

		var data map[string]any
		err := req.JSON(&data)
		must.NotNil(err)
		must.Nil(data)
	})
}

func TestFormValue(t *testing.T) {

	t.Run("can supply default value", func(t *testing.T) {
		must := must.New(t)
		req := Req(httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")))
		result := req.FormValue("idontknow", "1")
		must.Equal("1", result)
	})
}

type Key string

const KeyID Key = "id"

func TestLocals(t *testing.T) {
	must := must.New(t)
	req := Req(httptest.NewRequest(http.MethodGet, "/", nil))
	req.Locals(string(KeyID), 1)

	id := req.Locals(string(KeyID))
	must.Equal(1, id)

	slug := req.Locals("slug")
	must.Nil(slug)
}

func TestIsLocalhost(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	result := Req(req).IsLocalhost()
	must.False(result)

	req.Host = "localhost"
	result = Req(req).IsLocalhost()
	must.True(result)
}
