package ngamux

import "net/http"

// MiddlewareFunc describes a middleware function used by Ngamux.
// A middleware receives the next http.HandlerFunc and returns a new
// http.HandlerFunc that wraps the original handler. Use middlewares to
// implement cross-cutting concerns such as logging, authentication,
// CORS, request timeouts, and panic recovery.
type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

// ToHandler converts an http.HandlerFunc into an http.Handler. This is a
// convenience adapter useful when a function handler needs to be passed to
// APIs that require the http.Handler interface.
func ToHandler(h http.HandlerFunc) http.Handler {
	return http.Handler(h)
}

// ToHandlerFunc converts an http.Handler into an http.HandlerFunc by
// returning a function that forwards ServeHTTP to the provided handler.
// This is helpful when you have an http.Handler value and need the
// function form (http.HandlerFunc) for middleware composition or registration.
func ToHandlerFunc(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(h.ServeHTTP)
}
