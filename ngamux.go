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

// KeyContext describe key type for ngamux context
type KeyContext int

const (
	// KeyContextParams is key context for url params
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

// Use register global middleware
func (mux *Ngamux) Use(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

// Config returns registered config (read only)
func (mux Ngamux) Config() Config {
	return *mux.config
}

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

// Get register route for a url with Get request method
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

// Head register route for a url with Head request method
func (mux *Ngamux) Head(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(
		http.MethodHead,
		url,
		func(w http.ResponseWriter, r *http.Request) { handler(headResponseWriter{w}, r) },
		middlewares...,
	)
}

// Post register route for a url with Post request method
func (mux *Ngamux) Post(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodPost, url, handler, middlewares...)
}

// Patch register route for a url with Patch request method
func (mux *Ngamux) Patch(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodPatch, url, handler, middlewares...)
}

// Put register route for a url with Put request method
func (mux *Ngamux) Put(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodPut, url, handler, middlewares...)
}

// Delete register route for a url with Delete request method
func (mux *Ngamux) Delete(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc(http.MethodDelete, url, handler, middlewares...)
}

func (mux *Ngamux) All(url string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	mux.HandleFunc("ALL", url, handler, middlewares...)
}

// With register middlewares and returns router
func (mux *Ngamux) With(middlewares ...MiddlewareFunc) *Ngamux {
	group := mux.Group(mux.path)
	group.Use(middlewares...)
	return group
}
