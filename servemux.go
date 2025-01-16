package ngamux

import (
	"fmt"
	"net/http"
	"slices"
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
	_, pattern := h.mux.Handler(r)
	if pattern == "" || (pattern == "GET /" && r.URL.Path != "/") {
		WithMiddlewares(h.middlewares...)(http.NotFound).ServeHTTP(w, r)
		return
	}
	h.mux.ServeHTTP(w, r)
}

func (h *HttpServeMux) HandleFunc(method, path string, handlerFunc http.HandlerFunc, middlewares ...MiddlewareFunc) {
	slices.Reverse(middlewares)
	middlewares = append(h.middlewares, middlewares...)
	if h.parent != nil {
		if path == "/" {
			path = ""
		}
		route := fmt.Sprintf("%s %s%s", method, h.path, path)
		middlewares := append(h.parent.middlewares, middlewares...)
		h.parent.mux.Handle(route, WithMiddlewares(middlewares...)(handlerFunc))
		return
	}

	route := fmt.Sprintf("%s %s", method, path)
	h.mux.HandleFunc(route, WithMiddlewares(h.middlewares...)(handlerFunc))
}

func (h *HttpServeMux) Group(path string) *HttpServeMux {
	res := &HttpServeMux{
		path,
		http.NewServeMux(),
		h,
		make([]MiddlewareFunc, 0),
	}
	return res
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

//	func (mux *HttpServeMux) getParent() Router {
//		return mux.parent
//	}
func (mux *HttpServeMux) getPath() string {
	return mux.path
}
