package ngamux

import (
	"fmt"
	"net/http"
)

type HttpServeMux struct {
	path        string
	mux         *http.ServeMux
	parent      *HttpServeMux
	middlewares []MiddlewareFunc
}

func NewHttpServeMux() *HttpServeMux {
	return &HttpServeMux{
		"",
		http.NewServeMux(),
		nil,
		make([]MiddlewareFunc, 0),
	}
}

func (h *HttpServeMux) Use(middlewares ...MiddlewareFunc) {
	h.middlewares = append(h.middlewares, middlewares...)
}

func (h HttpServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, pattern := h.mux.Handler(r)
	if pattern == "GET /" && r.URL.Path != "/" {
		WithMiddlewares(h.middlewares...)(http.NotFound).ServeHTTP(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}

func (h *HttpServeMux) HandleFunc(method, path string, handlerFunc http.HandlerFunc) {
	if h.parent != nil {
		if path == "/" {
			path = ""
		}
		route := fmt.Sprintf("%s %s%s", method, h.path, path)
		middlewares := make([]MiddlewareFunc, 0)
		middlewares = append(middlewares, h.parent.middlewares...)
		middlewares = append(middlewares, h.middlewares...)
		h.parent.mux.Handle(route, WithMiddlewares(middlewares...)(handlerFunc))
		return
	}

	route := fmt.Sprintf("%s %s", method, path)
	h.mux.HandleFunc(route, handlerFunc)
}

func (h *HttpServeMux) Group(path string) *HttpServeMux {
	res := &HttpServeMux{
		path,
		http.NewServeMux(),
		h,
		h.middlewares,
	}
	return res
}

func (h *HttpServeMux) Get(path string, handlerFunc http.HandlerFunc) {
	h.HandleFunc(http.MethodGet, path, handlerFunc)
}

func (h *HttpServeMux) Post(path string, handlerFunc http.HandlerFunc) {
	h.HandleFunc(http.MethodPost, path, handlerFunc)
}

func (h *HttpServeMux) Patch(path string, handlerFunc http.HandlerFunc) {
	h.HandleFunc(http.MethodPatch, path, handlerFunc)
}

func (h *HttpServeMux) Put(path string, handlerFunc http.HandlerFunc) {
	h.HandleFunc(http.MethodPut, path, handlerFunc)
}

func (h *HttpServeMux) Delete(path string, handlerFunc http.HandlerFunc) {
	h.HandleFunc(http.MethodDelete, path, handlerFunc)
}
