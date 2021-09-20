package ngamux

import "net/http"

type (
	MiddlewareFunc func(next Handler) Handler
	Handler        func(rw http.ResponseWriter, r *http.Request) error
)

func (h Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h(rw, r)
}
