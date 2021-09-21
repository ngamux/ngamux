package ngamux

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type (
	Map map[string]interface{}
)

// SETUP
func WithMiddlewares(middleware ...MiddlewareFunc) MiddlewareFunc {
	return func(next Handler) Handler {
		h := next
		for i := len(middleware) - 1; i >= 0; i-- {
			h = middleware[i](h)
		}
		return h
	}
}

// REQUEST
func GetParam(r *http.Request, key string) string {
	params := r.Context().Value(KeyContextParams).([][]string)
	for _, param := range params {
		if param[0] == key {
			return param[1]
		}
	}

	return ""
}

func GetQuery(r *http.Request, key string, fallback ...string) string {
	queries := r.URL.Query()
	query := queries.Get(key)
	if query == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return ""
	}
	return query
}

func GetFormValue(r *http.Request, key string, fallback ...string) string {
	value := r.PostFormValue(key)
	if value == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return ""
	}

	return value
}

func GetFormFile(r *http.Request, key string, maxFileSize ...int64) (*multipart.FileHeader, error) {
	var maxFileSizeParsed int64 = 10 << 20
	if len(maxFileSize) > 0 {
		maxFileSizeParsed = maxFileSize[0]
	}

	r.ParseMultipartForm(maxFileSizeParsed)
	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return header, nil
}

func GetJSON(r *http.Request, store interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		return err
	}

	return nil
}

func SetContextValue(r *http.Request, key interface{}, value interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	r = r.WithContext(ctx)
	return r
}

func GetContextValue(r *http.Request, key interface{}) interface{} {
	value := r.Context().Value(key)
	return value
}

// RESPONSE
func String(rw http.ResponseWriter, data string) error {
	rw.Header().Add("content-type", "text/plain")
	fmt.Fprintln(rw, data)
	return nil
}

func StringWithStatus(rw http.ResponseWriter, status int, data string) error {
	String(rw, data)
	rw.WriteHeader(status)
	return nil
}

func JSON(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Add("content-type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprint(rw, string(jsonData))
	return nil
}

func JSONWithStatus(rw http.ResponseWriter, status int, data interface{}) error {
	rw.WriteHeader(status)
	err := JSON(rw, data)
	if err != nil {
		return err
	}

	return nil
}
