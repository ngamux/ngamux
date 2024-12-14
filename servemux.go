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
	h.mux.ServeHTTP(w, r)
}

func (h *HttpServeMux) HandlerFunc(method, path string, handlerFunc Handler) {
	if h.parent != nil {
		if path == "/" {
			path = ""
		}
		route := fmt.Sprintf("%s %s%s", method, h.path, path)
		middlewares := make([]MiddlewareFunc, 0)
		middlewares = append(middlewares, h.parent.middlewares...)
		middlewares = append(middlewares, h.middlewares...)
		h.parent.mux.HandleFunc(route, WithMiddlewares(middlewares...)(handlerFunc).ServeHTTP)
		return
	}

	route := fmt.Sprintf("%s %s", method, path)
	h.mux.HandleFunc(route, WithMiddlewares(h.middlewares...)(handlerFunc).ServeHTTP)
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

func (h *HttpServeMux) Get(path string, handlerFunc Handler) {
	h.HandlerFunc(http.MethodGet, path, handlerFunc)
}

func (h *HttpServeMux) Post(path string, handlerFunc Handler) {
	h.HandlerFunc(http.MethodPost, path, handlerFunc)
}

func (h *HttpServeMux) Patch(path string, handlerFunc Handler) {
	h.HandlerFunc(http.MethodPatch, path, handlerFunc)
}

func (h *HttpServeMux) Put(path string, handlerFunc Handler) {
	h.HandlerFunc(http.MethodPut, path, handlerFunc)
}

func (h *HttpServeMux) Delete(path string, handlerFunc Handler) {
	h.HandlerFunc(http.MethodDelete, path, handlerFunc)
}
