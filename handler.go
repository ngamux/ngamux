package ngamux

import "net/http"

type (
	// MiddlewareFunc describe middleware function
	MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc
)

func ToHandler(h http.HandlerFunc) http.Handler {
	return http.Handler(h)
}

func ToHandlerFunc(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(h.ServeHTTP)
}
