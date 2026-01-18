package ngamux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"maps"
	"net/http"
)

// Response define single response manager
type Response struct {
	http.ResponseWriter
	status int
}

// Res needs http.ResponseWriter and returns *Response object
func Res(rw http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: rw,
	}
}

func (r Response) statusSafe() int {
	status := r.status
	if status == 0 {
		status = http.StatusOK
	}

	return status
}

// Status write status code
func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}

// Text writes text/plain data with simple string as response body
func (r *Response) Text(data string) {
	r.WriteHeader(r.statusSafe())
	r.Header().Add("Content-Type", "text/plain")
	_, _ = fmt.Fprintln(r, data)
}

// JSON write application/json data with json encoded string as response body
func (r *Response) JSON(data any) {
	buf := poolByte.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		poolByte.Put(buf)
	}()

	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}

	maps.Copy(r.Header(), headerContentTypeJSON)
	r.WriteHeader(r.statusSafe())
	_, _ = r.Write(buf.Bytes())
}

// HTML write text/html data with HTML string as response body
func (r *Response) HTML(path string, data any) {
	r.WriteHeader(r.statusSafe())
	r.Header().Set("Content-Type", "text/html; charset=utf-8")

	temp, err := template.ParseFiles(path)
	if err != nil {
		return
	}

	err = temp.Execute(r, data)
	if err != nil {
		return
	}
}
