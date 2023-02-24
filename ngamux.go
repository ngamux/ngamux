// Ngamux is simple HTTP router for Go
// Github Repository: https://github.com/ngamux/ngamux
// Examples: https://github.com/ngamux/ngamux-example

// Package ngamux is simple HTTP router for Go that compatible with net/http,
// the standard library to serve HTTP. Designed to make everything goes
// in simple way.
package ngamux

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

// KeyContext describe key type for ngamux context
type KeyContext int

const (
	// KeyContextParams is key context for url params
	KeyContextParams KeyContext = 1 << iota
)

type (
	// Ngamux describe structure of ngamux object
	Ngamux struct {
		parent            *Ngamux
		path              string
		routes            routeMap
		routesParam       routeMap
		config            Config
		regexpParamFinded *regexp.Regexp
		middlewares       []MiddlewareFunc
	}
)

var (
	paramsFinder       = regexp.MustCompile("(:[a-zA-Z]+[0-9a-zA-Z]*)")
	globalErrorHandler = func(rw http.ResponseWriter, r *http.Request) error {
		err := Req(r).Locals("error").(error)
		if errors.Is(err, ErrorNotFound) {
			rw.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, ErrorMethodNotAllowed) {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}

		fmt.Fprintln(rw, err)
		return nil
	}

	allMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	// ErrorNotFound is errors object when searching failure
	ErrorNotFound = errors.New("not found")

	// ErrorMethodNotAllowed is errors object when there access to invalid method
	ErrorMethodNotAllowed = errors.New("method not allowed")
)

// New returns new ngamux object
func New(opts ...func(*Config)) *Ngamux {
	config := NewConfig()
	for _, opt := range opts {
		opt(&config)
	}

	routesMap := routeMap{}
	routesParamMap := routeMap{}
	router := &Ngamux{
		routes:            routesMap,
		routesParam:       routesParamMap,
		config:            config,
		regexpParamFinded: paramsFinder,
	}

	return router
}

// ServeHTTP run ngamux router matcher
func (mux *Ngamux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, r := mux.getRoute(r)
	err := route.Handler(rw, r)
	if err != nil {
		rw.WriteHeader(500)
		_, _ = rw.Write([]byte(err.Error()))
	}
}

// Use register global middleware
func (mux *Ngamux) Use(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
	mux.config.GlobalErrorHandler = WithMiddlewares(mux.middlewares...)(mux.config.GlobalErrorHandler)
}

// Get register route for a url with Get request method
func (mux *Ngamux) Get(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodGet, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodGet, handler, mux.middlewares...))
}

// Post register route for a url with Post request method
func (mux *Ngamux) Post(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPost, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPost, handler, mux.middlewares...))
}

// Patch register route for a url with Patch request method
func (mux *Ngamux) Patch(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPatch, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPatch, handler, mux.middlewares...))
}

// Put register route for a url with Put request method
func (mux *Ngamux) Put(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPut, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPut, handler, mux.middlewares...))
}

// Delete register route for a url with Delete request method
func (mux *Ngamux) Delete(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodDelete, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodDelete, handler, mux.middlewares...))
}

// All register route for a url with any request method
func (mux *Ngamux) All(url string, handler Handler) {
	for _, method := range allMethods {
		if mux.parent != nil {
			mux.addRouteFromGroup(buildRoute(url, method, handler))
			return
		}

		mux.addRoute(buildRoute(url, method, handler, mux.middlewares...))
	}
}

// With register middlewares and returns router
func (mux *Ngamux) With(middlewares ...MiddlewareFunc) *Ngamux {
	group := &Ngamux{
		parent:      mux,
		path:        mux.path,
		middlewares: middlewares,
	}
	return group
}
