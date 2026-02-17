// Ngamux is simple HTTP router for Go
// Github Repository: https://github.com/ngamux/ngamux
// Examples: https://github.com/ngamux/ngamux-example

// Package ngamux is simple HTTP router for Go that compatible with net/http,
// the standard library to serve HTTP. Designed to make everything goes
// in simple way.
package ngamux

import (
	"net/http"
	gopath "path"
	"slices"

	"github.com/ngamux/ngamux/mapping"
)

// KeyContext describes keys used when storing values in request contexts.
// It is defined as an int-based type to avoid collisions with other
// context keys from external packages.
type KeyContext int

const (
	// KeyContextParams is the context key under which route parameters
	// (extracted during routing) are stored. The value associated with
	// this key is [][]string where each element is a two-item slice
	// [name, value]. Handlers can read parameters using Req(r).Params(name).
	KeyContextParams KeyContext = 1 << iota
)

var (
	headerContentTypeJSON = http.Header{
		"Content-Type": []string{"application/json"},
	}
	headerContentTypeText = http.Header{
		"Content-Type": []string{"text/plain"},
	}
)

type Ngamux struct {
	root        mapping.Mapping[string, *Node]
	middlewares []MiddlewareFunc
	config      *Config
	path        string
	parent      *Ngamux
}

// New constructs a new Ngamux router. Optionally pass functional
// configuration options produced by functions in options.go, for example
// ngamux.New(ngamux.WithTrailingSlash(), ngamux.WithLogLevel(slog.LevelInfo)).
//
// The returned *Ngamux is ready to register routes with methods like
// Get, Post, or HandleFunc and can be used directly as an http.Handler.
func New(opts ...func(*Config)) *Ngamux {
	config := NewConfig()
	for _, opt := range opts {
		opt(&config)
	}
	return &Ngamux{
		root:        mapping.New[string, *Node](),
		middlewares: make([]MiddlewareFunc, 0),
		config:      &config,
	}
}

// Use registers middleware functions that will be applied to all routes
// registered on this Ngamux instance (but not automatically applied to
// parent groups). Middlewares are applied in the order they are added;
// the execution order wraps handlers such that the most recently added
// middleware runs first when a request is served.
func (mux *Ngamux) Use(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

// Config returns a copy of the active configuration for this router.
// The returned Config is a value copy so callers cannot mutate the
// router's configuration by accident.
func (mux Ngamux) Config() Config {
	return *mux.config
}

// HandleFunc registers a handler function for a given HTTP method and
// path. Middlewares passed here are combined with middlewares registered
// on this router and any parent groups. If this Ngamux is nested (created
// via Group), the final path and middleware chain are composed by walking
// the parent chain and joining path prefixes.
func (mux *Ngamux) HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	middlewares = append(mux.middlewares, middlewares...)
	if mux.parent == nil {
		mux.Handle(method+" "+path, WithMiddlewares(middlewares...)(handler))
		return
	}

	parent := mux.parent
	path = gopath.Join(mux.path, path)
	for parent != nil {
		path = gopath.Join(parent.path, path)
		middlewares = append(parent.middlewares, middlewares...)
		if parent.parent == nil {
			break
		}

		parent = parent.parent
	}

	parent.Handle(method+" "+path, WithMiddlewares(middlewares...)(handler))
}

// Get registers a handler for GET requests on the provided URL.
func (mux *Ngamux) Get(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodGet, url, handler, middlewares...)
}

type headResponseWriter struct {
	http.ResponseWriter
}

func (headResponseWriter) Write(in []byte) (int, error) {
	return 0, nil
}

// Head registers a handler for HEAD requests. It adapts the provided
// handler by wrapping the ResponseWriter so that the body is suppressed
// while headers and status codes are preserved.
func (mux *Ngamux) Head(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(
		http.MethodHead,
		url,
		func(w http.ResponseWriter, r *http.Request) { handler(headResponseWriter{w}, r) },
		middlewares...,
	)
}

// Post registers a handler for POST requests on the provided URL.
func (mux *Ngamux) Post(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodPost, url, handler, middlewares...)
}

// Patch registers a handler for PATCH requests on the provided URL.
func (mux *Ngamux) Patch(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodPatch, url, handler, middlewares...)
}

// Put registers a handler for PUT requests on the provided URL.
func (mux *Ngamux) Put(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodPut, url, handler, middlewares...)
}

// Delete registers a handler for DELETE requests on the provided URL.
func (mux *Ngamux) Delete(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodDelete, url, handler, middlewares...)
}

// All registers a handler that accepts requests of any HTTP method for the
// given URL. This is useful for endpoints that intentionally handle
// multiple methods in a single function.
func (mux *Ngamux) All(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc("ALL", url, handler, middlewares...)
}

// With creates a new sub-router (group) based on the current router's
// path and registers the provided middlewares on that group. This is a
// convenience for applying a short-lived middleware chain to a set of
// routes.
func (mux *Ngamux) With(middlewares ...MiddlewareFunc) *Ngamux {
	group := mux.Group(mux.path)
	group.Use(middlewares...)
	return group
}
