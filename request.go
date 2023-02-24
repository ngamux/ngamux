package ngamux

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
)

// Request define single request manager
type Request struct {
	*http.Request
}

// Req needs *http.Request and returns *Request object
func Req(r *http.Request) *Request {
	return &Request{r}
}

// Params returns parameter from url using a key
func (r Request) Params(key string) string {
	params := r.Context().Value(KeyContextParams).([][]string)
	for _, param := range params {
		if param[0] == key {
			return param[1]
		}
	}

	return ""
}

// Query returns data from query params using a key
func (r Request) Query(key string, fallback ...string) string {
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

// FormValue returns data from form using a key
func (r Request) FormValue(key string, fallback ...string) string {
	value := r.PostFormValue(key)
	if value == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return ""
	}

	return value
}

// FormFile returns file from form using a key
func (r Request) FormFile(key string, maxFileSize ...int64) (*multipart.FileHeader, error) {
	var maxFileSizeParsed int64 = 10 << 20
	if len(maxFileSize) > 0 {
		maxFileSizeParsed = maxFileSize[0]
	}

	err := r.ParseMultipartForm(maxFileSizeParsed)
	if err != nil {
		return nil, err
	}

	file, header, err := r.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return header, nil
}

// JSON get json data from request body and store to variable reference
func (r Request) JSON(store any) error {
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		return err
	}

	return nil
}

// Locals needs key and optional value
// It returns any if only key and no value given
// It insert value to context if key and value is given
func (r *Request) Locals(key any, value ...any) any {
	if len(value) <= 0 {
		return r.Context().Value(key)
	}
	r.Request = r.WithContext(context.WithValue(r.Context(), key, value[0]))
	return nil
}
