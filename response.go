package ngamux

import (
	"encoding/json"
	"fmt"
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

// String write string data to response body
func (r *Response) String(data string) error {
	r.WriteHeader(r.statusSafe())
	r.Header().Add("content-type", "text/plain")
	_, err := fmt.Fprintln(r, data)

	return err
}

// JSON write JSON data to response
func (r *Response) JSON(data any) error {
	r.WriteHeader(r.statusSafe())
	r.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(r).Encode(data); err != nil {
		return err
	}

	return nil
}
