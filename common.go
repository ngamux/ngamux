package ngamux

import "net/http"

type (
	// Map is key value type to store any data
	Map map[string]any
)

// SETUP

// WithMiddlewares returns single middleware from multiple middleware
func WithMiddlewares(middleware ...MiddlewareFunc) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		h := next
		if len(middleware) <= 0 {
			return h
		}

		for i := len(middleware) - 1; i >= 0; i-- {
			if middleware[i] == nil {
				continue
			}
			h = middleware[i](h)
		}
		return h
	}
}
