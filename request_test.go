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

func TestGetParam(t *testing.T) {
	must := must.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), KeyContextParams, [][]string{{"id", "1"}}))
	result := Req(req).GetParam("id")
	must.Equal("1", result)

	result = Req(req).GetParam("slug")
	must.Equal("", result)
}

func TestGetQuery(t *testing.T) {
	must := must.New(t)
	req := Req(httptest.NewRequest(http.MethodGet, "/?id=1", nil))
	result := req.GetQuery("id")
	must.Equal("1", result)

	result = req.GetQuery("slug", "undefined")
	must.Equal("undefined", result)

	result = req.GetQuery("slug")
	must.Equal("", result)
}

func TestGetJSON(t *testing.T) {
	must := must.New(t)
	input := strings.NewReader(`{"id": 1}`)
	req := Req(httptest.NewRequest(http.MethodGet, "/?id=1", input))

	var data map[string]any
	err := req.GetJSON(&data)
	must.Nil(err)
	must.NotNil(data["id"])

	id, ok := data["id"].(float64)
	must.True(ok)
	must.Equal(float64(1), id)
}

func TestGetFormValue(t *testing.T) {

	t.Run("can supply default value", func(t *testing.T) {
		must := must.New(t)
		req := Req(httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")))
		result := req.GetFormValue("idontknow", "1")
		must.Equal("1", result)
	})
}

func TestSetContextValue(t *testing.T) {
	must := must.New(t)
	req := Req(httptest.NewRequest(http.MethodGet, "/", nil))
	req = req.SetContextValue("id", 1)

	id := req.Context().Value("id")
	must.Equal(1, id)

	slug := req.Context().Value("slug")
	must.Nil(slug)
}

type Key string

const KeyID Key = "id"

func TestGetContextValue(t *testing.T) {
	must := must.New(t)
	req := Req(httptest.NewRequest(http.MethodGet, "/", nil))
	req = req.SetContextValue(KeyID, 1)

	id := req.GetContextValue(KeyID)
	must.Equal(1, id)

	slug := req.GetContextValue("slug")
	must.Nil(slug)
}
