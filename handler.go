package ngamux

import "net/http"

type (
	// MiddlewareFunc describe middleware function
	MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

	// Handler describe function handler
	handler func(rw http.ResponseWriter, r *http.Request)
)

// ServeHTTP same as original Handler but for built in HTTP HandlerFunc
func (h handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h(rw, r)
}

func ToHandler(h http.HandlerFunc) http.Handler {
	return http.Handler(h)
}

func ToHandlerFunc(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(h.ServeHTTP)
}
