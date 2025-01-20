package ngamux

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestGetIpAddress(t *testing.T) {

	must := must.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Inject IP Local
	req.RemoteAddr = "127.0.0.1"
	result := Req(req).GetIPAdress()
	must.Equal("127.0.0.1", result)
}

type queryStruct struct {
	Foo string `query:"foo"`
	Bar string `query:"bar"`
}

type queryStructWithNested struct {
	FirstName string  `query:"first_name"`
	LastName  string  `query:"last_name"`
	Age       int     `query:"age"`
	Wallet    float64 `query:"wallet"`
	IsHuman   bool    `query:"is_human"`
	Address   struct {
		City    string `query:"city"`
		Country string `query:"country"`
	} `query:"address"`
}

type queryStructNotValid struct {
	FieldNotValid interface{} `query:"field_not_valid"` // <- dont use interface{} to data types
}

func TestQueriesParser(t *testing.T) {
	t.Run("parse query url to struct", func(t *testing.T) {
		must := must.New(t)
		req := httptest.NewRequest(http.MethodGet, "/q?foo=foo&bar=bar", nil)

		var qs queryStruct
		err := Req(req).QueriesParser(&qs)
		must.Nil(err)

		must.Equal(qs.Foo, "foo")
		must.Equal(qs.Bar, "bar")
	})

	t.Run("parse query url to nested struct", func(t *testing.T) {
		must := must.New(t)
		req := httptest.NewRequest(http.MethodGet, "/q?age=19&is_human=true&wallet=69.8&first_name=farda&last_name=nurfatika&address=city%3Dsemarang%26country%3Dindonesia", nil)

		var qswn queryStructWithNested
		err := Req(req).QueriesParser(&qswn)
		must.Nil(err)

		must.Equal(qswn.FirstName, "farda")
		must.Equal(qswn.LastName, "nurfatika")
		must.Equal(qswn.Age, 19)
		must.Equal(qswn.IsHuman, true)
		must.Equal(qswn.Wallet, 69.8)
		must.Equal(qswn.Address.City, "semarang")
		must.Equal(qswn.Address.Country, "indonesia")
	})

	t.Run("parse query to struct error", func(t *testing.T) {
		must := must.New(t)
		req := httptest.NewRequest(http.MethodGet, "/q?field_not_valid=xixixixi", nil)

		var qs queryStructNotValid
		err := Req(req).QueriesParser(&qs)
		must.NotNil(err)

		must.Equal(qs.FieldNotValid, nil)
	})
}

func TestFormFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Name())
	must.Nil(t, err)

	_, err = io.Copy(part, file)
	must.Nil(t, err)

	err = writer.Close()
	must.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	headers, err := Req(req).FormFile("file", 1)
	must.Nil(t, err)
	must.True(t, strings.HasSuffix(file.Name(), headers.Filename))
}
