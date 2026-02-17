package ngamux

import (
	"fmt"
	"net/http"
	gopath "path"
	"slices"

	"github.com/ngamux/ngamux/json"
)

// HttpServeMux is a lightweight wrapper around the standard library's
// http.ServeMux that adds support for middleware stacks and simple
// subgrouping (nested groups with shared path prefixes and middlewares).
//
// This type is intentionally simpler than Ngamux and exists to provide a
// convenience adapter for users who prefer the standard ServeMux behavior
// but still want middleware support.
type HttpServeMux struct {
	path        string
	mux         *http.ServeMux
	parent      *HttpServeMux
	middlewares []MiddlewareFunc
	config      *Config
}

// NewHttpServeMux constructs a new HttpServeMux. Optionally a Config can be
// provided to override defaults (JSON marshalling, logging level, etc.).
func NewHttpServeMux(cfg ...*Config) *HttpServeMux {
	if len(cfg) <= 0 {
		c := NewConfig()
		cfg = append(cfg, &c)
	}

	json.Configure(cfg[0].JSONMarshal, cfg[0].JSONUnmarshal)

	return &HttpServeMux{
		"",
		http.NewServeMux(),
		nil,
		make([]MiddlewareFunc, 0),
		cfg[0],
	}
}

// Use registers one or more middleware functions on this ServeMux group.
// Middlewares registered on a group are applied to routes added to that
// group (and to sub-groups when you register them via Group()).
func (h *HttpServeMux) Use(middlewares ...MiddlewareFunc) {
	h.middlewares = append(h.middlewares, middlewares...)
}

// ServeHTTP implements http.Handler. It delegates to the underlying
// http.ServeMux but first checks whether a matching pattern exists. If no
// route matches, registered middlewares will be applied to the
// http.NotFound handler.
func (h HttpServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pattern := h.mux.Handler(r)
	if pattern == "" || (pattern == "GET /" && r.URL.Path != "/") {
		WithMiddlewares(h.middlewares...)(http.NotFound).ServeHTTP(w, r)
		return
	}
	h.mux.ServeHTTP(w, r)
}

// HandleFunc registers a handler function for a method and path. If this
// group has a parent, the effective path and middleware stack are built up
// from parent groups so that nested groups inherit path prefixes and
// middlewares.
func (h *HttpServeMux) HandleFunc(method, path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	if h.parent == nil {
		route := fmt.Sprintf("%s %s", method, path)
		middlewares = append(h.middlewares, middlewares...)
		h.mux.HandleFunc(route, WithMiddlewares(middlewares...)(handlerFunc))
		return
	}

	parent := h
	paths := []string{path}
	for {
		paths = append([]string{parent.path}, paths...)
		middlewares = append(parent.middlewares, middlewares...)

		if parent.parent == nil {
			break
		}
		parent = parent.parent
	}
	route := fmt.Sprintf("%s %s", method, gopath.Join(paths...))
	parent.mux.HandleFunc(route, WithMiddlewares(middlewares...)(handlerFunc))
}

// Group creates a nested HttpServeMux group with a path prefix. The new
// group inherits configuration from the parent but can register its own
// middlewares.
func (h *HttpServeMux) Group(path string) *HttpServeMux {
	res := &HttpServeMux{
		path,
		http.NewServeMux(),
		h,
		make([]MiddlewareFunc, 0),
		h.config,
	}
	return res
}

// GroupFunc creates a group and invokes the provided router function so
// the caller can register routes on the returned group inline.
func (h *HttpServeMux) GroupFunc(path string, router func(mux *HttpServeMux)) {
	group := h.Group(path)
	router(group)
}

func (h *HttpServeMux) Get(path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	h.HandleFunc(http.MethodGet, path, handlerFunc, middlewares...)
}

func (h *HttpServeMux) Post(path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	h.HandleFunc(http.MethodPost, path, handlerFunc, middlewares...)
}

func (h *HttpServeMux) Patch(path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	h.HandleFunc(http.MethodPatch, path, handlerFunc, middlewares...)
}

func (h *HttpServeMux) Put(path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	h.HandleFunc(http.MethodPut, path, handlerFunc, middlewares...)
}

func (h *HttpServeMux) Delete(path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	h.HandleFunc(http.MethodDelete, path, handlerFunc, middlewares...)
}

func (h *HttpServeMux) All(path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	h.HandleFunc("", path, handlerFunc, middlewares...)
}
