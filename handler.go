package ngamux

import "net/http"

type (
	// MiddlewareFunc describe middleware function
	MiddlewareFunc func(next Handler) Handler

	// Handler describe function handler
	Handler func(rw http.ResponseWriter, r *http.Request) error
)

// ServeHTTP same as original Handler but for built in HTTP HandlerFunc
func (h Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h(rw, r)
}
