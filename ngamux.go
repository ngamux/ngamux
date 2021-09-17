package ngamux

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"regexp"
)

type KeyContext int

const (
	KeyContextParams KeyContext = 1 << iota
)

type (
	MiddlewareFunc func(next HandlerFunc) HandlerFunc
	HandlerFunc    func(rw http.ResponseWriter, r *http.Request) error
)

type Config struct {
	RemoveTrailingSlash bool
	NotFoundHandler     HandlerFunc
}

type Ngamux struct {
	parent            *Ngamux
	path              string
	routes            map[string]map[string]Route
	routesParam       map[string]map[string]Route
	config            Config
	regexpParamFinded *regexp.Regexp
	middlewares       []MiddlewareFunc
}

var _ http.Handler = &Ngamux{}

type Route struct {
	Path       string
	Handler    HandlerFunc
	Params     [][]string
	UrlMathcer *regexp.Regexp
}

var handlerNotFound = func(rw http.ResponseWriter, r *http.Request) error {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(rw, "404 page not found")
	return nil
}

func buildRoute(url string, handler HandlerFunc, middlewares ...MiddlewareFunc) Route {
	handler = WithMiddlewares(middlewares...)(handler)

	return Route{
		Path:    url,
		Handler: handler,
	}
}

func makeConfig(configs ...Config) Config {
	config := Config{
		RemoveTrailingSlash: true,
	}
	if len(configs) > 0 {
		config = configs[0]
	}
	if config.NotFoundHandler == nil {
		config.NotFoundHandler = handlerNotFound
	}

	return config
}

func NewNgamux(configs ...Config) *Ngamux {
	config := makeConfig(configs...)

	routesMap := map[string]map[string]Route{
		http.MethodGet:     {},
		http.MethodPost:    {},
		http.MethodPatch:   {},
		http.MethodPut:     {},
		http.MethodDelete:  {},
		http.MethodOptions: {},
		http.MethodConnect: {},
		http.MethodHead:    {},
		http.MethodTrace:   {},
	}

	paramsFinder, err := regexp.Compile("(:[a-zA-Z][0-9a-zA-Z]+)")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	router := &Ngamux{
		routes:            routesMap,
		routesParam:       make(map[string]map[string]Route),
		config:            config,
		regexpParamFinded: paramsFinder,
	}

	for key, val := range router.routes {

		var row = make(map[string]Route)
		for key2, route := range val {
			row[key2] = route
		}

		router.routesParam[key] = row
	}

	return router
}

func (mux *Ngamux) addRoute(method string, route Route) {
	var (
		err            error
		pathWithParams string
		subMatchs      = mux.regexpParamFinded.FindAllStringSubmatch(route.Path, -1)
	)

	// check if route doesn't have url param
	if len(subMatchs) == 0 {
		mux.routes[method][route.Path] = route
		return
	}

	route.Params = [][]string{}
	for _, val := range subMatchs {
		route.Params = append(route.Params, []string{val[0][1:]})
	}

	pathWithParams = mux.regexpParamFinded.ReplaceAllString(route.Path, "([0-9a-zA-Z]+)")
	route.Path = pathWithParams

	route.UrlMathcer, err = regexp.Compile("^" + pathWithParams + "$")
	if err != nil {
		log.Fatal(err)
		return
	}

	mux.routesParam[method][route.Path] = route
}

func (mux *Ngamux) addRouteFromGroup(method string, route Route) {
	url := path.Join(mux.path, route.Path)
	middlewares := mux.middlewares
	middlewares = append(middlewares, mux.parent.middlewares...)
	mux.parent.addRoute(method, buildRoute(url, route.Handler, middlewares...))
}

func (mux *Ngamux) getRoute(method string, path string) Route {
	foundRoute, ok := mux.routes[method][path]
	if !ok {
		for url, route := range mux.routesParam[method] {

			if route.UrlMathcer.MatchString(path) {
				foundParams := route.UrlMathcer.FindAllStringSubmatch(path, -1)
				params := make([][]string, len(route.Params))
				copy(params, route.Params)
				for i := range params {
					params[i] = append(params[i], foundParams[0][i+1])
				}
				route.Params = params
				foundRoute = route
				break
			}

			if url == path {
				foundRoute = route
				break
			}
		}
	}

	if foundRoute.Handler == nil {
		foundRoute.Handler = mux.config.NotFoundHandler
	}

	return foundRoute
}

func (mux *Ngamux) Group(url string, middlewares ...MiddlewareFunc) *Ngamux {
	group := &Ngamux{
		parent:      mux,
		path:        url,
		middlewares: middlewares,
	}
	return group
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

func (mux *Ngamux) Get(url string, handler HandlerFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodGet, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodGet, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Post(url string, handler HandlerFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodPost, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodPost, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Patch(url string, handler HandlerFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodPatch, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodPatch, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Put(url string, handler HandlerFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodPut, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodPut, buildRoute(url, handler, mux.middlewares...))
}

func (mux *Ngamux) Delete(url string, handler HandlerFunc) {
	if mux.parent != nil {
		mux.addRouteFromGroup(http.MethodDelete, buildRoute(url, handler))
		return
	}
	mux.addRoute(http.MethodDelete, buildRoute(url, handler, mux.middlewares...))
}
