package ngamux

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
)

type Request struct {
	*http.Request
}

func Req(r *http.Request) *Request {
	return &Request{r}
}

// GetParam returns parameter from url using a key
func (r Request) GetParam(key string) string {
	params := r.Context().Value(KeyContextParams).([][]string)
	for _, param := range params {
		if param[0] == key {
			return param[1]
		}
	}

	return ""
}

// GetQuery returns data from query params using a key
func (r Request) GetQuery(key string, fallback ...string) string {
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
func (r Request) GetFormValue(key string, fallback ...string) string {
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
func (r Request) GetFormFile(key string, maxFileSize ...int64) (*multipart.FileHeader, error) {
	var maxFileSizeParsed int64 = 10 << 20
	if len(maxFileSize) > 0 {
		maxFileSizeParsed = maxFileSize[0]
	}

	err := r.ParseMultipartForm(maxFileSizeParsed)
	if err != nil {
		return nil, err
	}

	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return header, nil
}

// GetJSON get json data from requst body and store to variable reference
func (r Request) GetJSON(store any) error {
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		return err
	}

	return nil
}

// SetContextValue returns http request object with new context that contains value
func (r *Request) SetContextValue(key any, value any) *Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	r = &Request{r.WithContext(ctx)}
	return r
}

// GetContextValue returns value from http request context
func (r Request) GetContextValue(key any) any {
	value := r.Context().Value(key)
	return value
}
