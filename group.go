package ngamux

import (
	"net/http"
	"path"
)

type group struct {
	parent      *Ngamux
	path        string
	middlewares []http.HandlerFunc
}

func (mux group) buildRoute(url string, handler ...http.HandlerFunc) Route {
	url = path.Join(mux.path, url)
	middlewares := append(mux.middlewares, handler...)
	handler = middlewares
	return Route{
		Path:     url,
		Handlers: handler,
	}
}

func (mux *group) Get(url string, handler ...http.HandlerFunc) {
	route := mux.buildRoute(url, handler...)
	mux.parent.router.AddRoute(http.MethodGet, route)
}

func (mux *group) Post(url string, handler ...http.HandlerFunc) {
	route := mux.buildRoute(url, handler...)
	mux.parent.router.AddRoute(http.MethodPost, route)
}

func (mux *group) Patch(url string, handler ...http.HandlerFunc) {
	route := mux.buildRoute(url, handler...)
	mux.parent.router.AddRoute(http.MethodPatch, route)
}

func (mux *group) Put(url string, handler ...http.HandlerFunc) {
	route := mux.buildRoute(url, handler...)
	mux.parent.router.AddRoute(http.MethodPut, route)
}

func (mux *group) Delete(url string, handler ...http.HandlerFunc) {
	route := mux.buildRoute(url, handler...)
	mux.parent.router.AddRoute(http.MethodDelete, route)
}
