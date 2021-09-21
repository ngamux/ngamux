package ngamux

import (
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
	routesMap := routeMap{}
	routesParamMap := routeMap{}
	router := &Ngamux{
		routes:            routesMap,
		routesParam:       routesParamMap,
		config:            config,
		regexpParamFinded: paramsFinder,
	}

	return router
}

func (mux *Ngamux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, r := mux.getRoute(r)
	route.Handler(rw, r)
}

func (mux *Ngamux) Use(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
	mux.config.NotFoundHandler = WithMiddlewares(mux.middlewares...)(mux.config.NotFoundHandler)
}

func (mux *Ngamux) Get(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodGet, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodGet, handler, mux.middlewares...))
}

func (mux *Ngamux) Post(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPost, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPost, handler, mux.middlewares...))
}

func (mux *Ngamux) Patch(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPatch, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPatch, handler, mux.middlewares...))
}

func (mux *Ngamux) Put(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPut, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPut, handler, mux.middlewares...))
}

func (mux *Ngamux) Delete(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodDelete, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodDelete, handler, mux.middlewares...))
}

func (mux *Ngamux) All(url string, handler Handler) {
	if mux.parent != nil {
		for method := range mux.routes {
			mux.addRouteFromGroup(buildRoute(url, method, handler))
		}
		return
	}

	for method := range mux.routes {
		mux.addRoute(buildRoute(url, method, handler, mux.middlewares...))
	}
}
