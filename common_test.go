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

	var data map[string]interface{}
	err := GetJSON(req, &data)
	must.Nil(err)
	must.NotNil(data["id"])

	id, ok := data["id"].(float64)
	must.True(ok)
	must.Equal(float64(1), id)
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

func TestString(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := String(rw, "ok")
		must.Nil(err)
	})
	handler.ServeHTTP(rec, req)

	result := rec.Body.String()
	expected := "ok\n"
	must.Equal(expected, result)
}

func TestStringWithStatus(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := StringWithStatus(rw, http.StatusOK, "ok")
		must.Nil(err)
	})
	handler.ServeHTTP(rec, req)

	resultBody := rec.Body.String()
	expectedBody := "ok\n"
	must.Equal(expectedBody, resultBody)

	resultStatus := rec.Result().StatusCode
	expectedStatus := http.StatusOK
	must.Equal(expectedStatus, resultStatus)
}

func TestJSON(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		JSON(rw, Map{
			"id": 1,
		})
	})
	handler.ServeHTTP(rec, req)

	resultBody := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expectedBody := `{"id":1}`
	must.Equal(expectedBody, resultBody)
}

func TestJSONWithStatus(t *testing.T) {
	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		JSONWithStatus(rw, http.StatusOK, Map{
			"id": 1,
		})
	})
	handler.ServeHTTP(rec, req)

	resultBody := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expectedBody := `{"id":1}`
	must.Equal(expectedBody, resultBody)

	resultStatus := rec.Result().StatusCode
	expectedStatus := http.StatusOK
	must.Equal(expectedStatus, resultStatus)
}
