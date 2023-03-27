package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestRes(t *testing.T) {
	must := must.New(t)
	rec := httptest.NewRecorder()
	expected := rec
	result := Res(rec)
	must.NotNil(result)
	must.Equal(expected, result.ResponseWriter)
}

func TestResStatus(t *testing.T) {
	must := must.New(t)
	expected := http.StatusOK
	r := Res(httptest.NewRecorder())
	r = r.Status(expected)

	must.NotNil(r)
	must.Equal(expected, r.status)
}

func TestResText(t *testing.T) {
	must := must.New(t)
	expected := "ok"
	result := httptest.NewRecorder()
	r := Res(result)
	err := r.Text(expected)

	must.Nil(err)
	must.Equal(r.status, 0)
	must.Equal(expected, strings.ReplaceAll(result.Body.String(), "\n", ""))
}

func TestResJSON(t *testing.T) {
	must := must.New(t)
	expected := `{"id":1}`
	result := httptest.NewRecorder()
	r := Res(result)
	err := r.JSON(Map{
		"id": 1,
	})

	must.Nil(err)
	must.Equal(r.status, 0)
	must.Equal(expected, strings.ReplaceAll(result.Body.String(), "\n", ""))
}

func TestResHtml(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		must := must.New(t)
		expected := `<span>mantap</span>`
		result := httptest.NewRecorder()
		r := Res(result)
		err := r.HTML("./response_test.html", nil)

		must.Nil(err)
		must.Equal(r.status, 0)
		must.Equal(expected, strings.ReplaceAll(result.Body.String(), "\n", ""))
	})

	t.Run("negative", func(t *testing.T) {
		must := must.New(t)
		expected := ""
		result := httptest.NewRecorder()
		r := Res(result)
		err := r.HTML("./response.html", nil)

		must.NotNil(err)
		must.Equal(r.status, 0)
		must.Equal(expected, strings.ReplaceAll(result.Body.String(), "\n", ""))
	})
}
