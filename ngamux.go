package ngamux

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

type KeyContext int

const (
	KeyContextParams KeyContext = 1 << iota
)

type (
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
	_               http.Handler = &Ngamux{}
	_               http.Handler = Handler(func(rw http.ResponseWriter, r *http.Request) error { return nil })
	paramsFinder                 = regexp.MustCompile("(:[a-zA-Z][0-9a-zA-Z]+)")
	handlerNotFound              = func(rw http.ResponseWriter, r *http.Request) error {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(rw, "404 page not found")
		return nil
	}
)

func NewNgamux(configs ...Config) *Ngamux {
	config := buildConfig(configs...)
	routesMap := buildRouteMap()
	routesParamMap := buildRouteMap()
	router := &Ngamux{
		routes:            routesMap,
		routesParam:       routesParamMap,
		config:            config,
		regexpParamFinded: paramsFinder,
	}

	return router
}

func (mux *Ngamux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if mux.config.RemoveTrailingSlash && len(url) > 1 && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	route := mux.getRoute(r.Method, url)
	if len(route.Params) > 0 {
		ctx := context.WithValue(r.Context(), KeyContextParams, route.Params)
		r = r.WithContext(ctx)
	}

	route.Handler(rw, r)
}

func (mux *Ngamux) Use(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

func (mux *Ngamux) Get(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodGet, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodGet, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Post(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodPost, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodPost, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Patch(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodPatch, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodPatch, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Put(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodPut, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodPut, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Delete(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodDelete, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodDelete, buildRoute(url, handler, mux.middlewares...))
}
