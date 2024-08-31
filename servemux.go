package ngamux

import (
	"fmt"
	"net/http"
)

type HttpServeMux struct {
	path   string
	mux    *http.ServeMux
	parent *HttpServeMux
}

func NewHttpServeMux() *HttpServeMux {
	return &HttpServeMux{
		"",
		http.NewServeMux(),
		nil,
	}
}

func (h HttpServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *HttpServeMux) HandlerFunc(method, path string, handlerFunc http.HandlerFunc) {
	if h.parent != nil {
		if path == "/" {
			path = ""
		}
		route := fmt.Sprintf("%s %s%s", method, h.path, path)
		h.parent.mux.HandleFunc(route, handlerFunc)
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
	}
	return res
}

func (h *HttpServeMux) Get(path string, handlerFunc http.HandlerFunc) {
	h.HandlerFunc(http.MethodGet, path, handlerFunc)
}

func (h *HttpServeMux) Post(path string, handlerFunc http.HandlerFunc) {
	h.HandlerFunc(http.MethodPost, path, handlerFunc)
}

func (h *HttpServeMux) Patch(path string, handlerFunc http.HandlerFunc) {
	h.HandlerFunc(http.MethodPatch, path, handlerFunc)
}

func (h *HttpServeMux) Put(path string, handlerFunc http.HandlerFunc) {
	h.HandlerFunc(http.MethodPut, path, handlerFunc)
}

func (h *HttpServeMux) Delete(path string, handlerFunc http.HandlerFunc) {
	h.HandlerFunc(http.MethodDelete, path, handlerFunc)
}
