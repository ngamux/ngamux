package ngamux

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// Response define single response manager
type Response struct {
	http.ResponseWriter
	status int
}

type readOnlyResponseWriter struct {
	http.ResponseWriter
}

func (r readOnlyResponseWriter) Write(data []byte) (int, error) {
	return 0, nil
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
func (r *Response) Text(data string) error {
	r.WriteHeader(r.statusSafe())
	r.Header().Add("content-type", "text/plain")
	_, err := fmt.Fprintln(r, data)

	return err
}

// JSON write application/json data with json encoded string as response body
func (r *Response) JSON(data any) error {
	r.WriteHeader(r.statusSafe())
	r.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(r).Encode(data); err != nil {
		return err
	}

	return nil
}

// HTML write text/html data with HTML string as response body
func (r *Response) HTML(path string, data any) error {
	r.WriteHeader(r.statusSafe())
	r.Header().Add("Content-Type", "text/html; charset=utf-8")

	temp, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	err = temp.Execute(r, data)
	if err != nil {
		return err
	}

	return nil
}
