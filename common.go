package ngamux

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type (
	// Map is key value type to store any data
	Map map[string]interface{}
)

// SETUP

// WithMiddlewares returns single middleware from multiple middleware
func WithMiddlewares(middleware ...MiddlewareFunc) MiddlewareFunc {
	return func(next Handler) Handler {
		h := next
		if len(middleware) <= 0 {
			return h
		}

		for i := len(middleware) - 1; i >= 0; i-- {
			if middleware[i] == nil {
				continue
			}
			h = middleware[i](h)
		}
		return h
	}
}

// REQUEST

// GetParam returns parameter from url using a key
func GetParam(r *http.Request, key string) string {
	params := r.Context().Value(KeyContextParams).([][]string)
	for _, param := range params {
		if param[0] == key {
			return param[1]
		}
	}

	return ""
}

// GetQuery returns data from query params using a key
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

// GetFormValue returns data from form using a key
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

// GetFormFile returns file from form using a key
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

// GetJSON get json data from requst body and store to variable reference
func GetJSON(r *http.Request, store interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		return err
	}

	return nil
}

// SetContextValue returns http request object with new context that contains value
func SetContextValue(r *http.Request, key interface{}, value interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	r = r.WithContext(ctx)
	return r
}

// GetContextValue returns value from http request context
func GetContextValue(r *http.Request, key interface{}) interface{} {
	value := r.Context().Value(key)
	return value
}

// RESPONSE

// String write string data to response body
func String(rw http.ResponseWriter, data string) error {
	rw.Header().Add("content-type", "text/plain")
	fmt.Fprintln(rw, data)
	return nil
}

// StringWithStatus write string data to response body with status code
func StringWithStatus(rw http.ResponseWriter, status int, data string) error {
	String(rw, data)
	rw.WriteHeader(status)
	return nil
}

// JSON write JSON data to response
func JSON(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return err
	}

	return nil
}

// JSONWithStatus write JSON data to response body with status code
func JSONWithStatus(rw http.ResponseWriter, status int, data interface{}) error {
	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(status)
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return err
	}

	return nil
}
