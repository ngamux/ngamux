package ngamux

import (
	"fmt"
	"net/http"
	gopath "path"
	"slices"
)

type HttpServeMux struct {
	path        string
	mux         *http.ServeMux
	parent      *HttpServeMux
	middlewares []MiddlewareFunc
	config      *Config
}

func NewHttpServeMux(cfg ...*Config) *HttpServeMux {
	if len(cfg) <= 0 {
		c := NewConfig()
		cfg = append(cfg, &c)
	}

	return &HttpServeMux{
		"",
		http.NewServeMux(),
		nil,
		make([]MiddlewareFunc, 0),
		cfg[0],
	}
}

func (h *HttpServeMux) Use(middlewares ...MiddlewareFunc) {
	h.middlewares = append(h.middlewares, middlewares...)
}

func (h HttpServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pattern := h.mux.Handler(r)
	if pattern == "" || (pattern == "GET /" && r.URL.Path != "/") {
		WithMiddlewares(h.middlewares...)(http.NotFound).ServeHTTP(w, r)
		return
	}
	h.mux.ServeHTTP(w, r)
}

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
