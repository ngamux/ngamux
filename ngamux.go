package ngamux

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

type KeyContext int

const (
	KeyContextParams KeyContext = 1 << iota
)

type (
	MiddlewareFunc func(next HandlerFunc) HandlerFunc
	HandlerFunc    func(rw http.ResponseWriter, r *http.Request) error

	Config struct {
		RemoveTrailingSlash bool
		NotFoundHandler     HandlerFunc
	}

	Ngamux struct {
		parent            *Ngamux
		path              string
		routes            routeMap
		routesParam       routeMap
		config            Config
		regexpParamFinded *regexp.Regexp
		middlewares       []MiddlewareFunc
	}

	Route struct {
		Path       string
		Handler    HandlerFunc
		Params     [][]string
		UrlMathcer *regexp.Regexp
	}

	routeMap map[string]map[string]Route
)

var _ http.Handler = &Ngamux{}

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

	routesMap := routeMap{
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

	routesParamMap := routeMap{
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

	paramsFinder := regexp.MustCompile("(:[a-zA-Z][0-9a-zA-Z]+)")
	router := &Ngamux{
		routes:            routesMap,
		routesParam:       routesParamMap,
		config:            config,
		regexpParamFinded: paramsFinder,
	}

	return router
}

func (mux *Ngamux) addRoute(method string, route Route) {
	var (
		err            error
		pathWithParams string
	)

	// check if route doesn't have url param
	if !strings.Contains(route.Path, ":") {
		mux.routes[method][route.Path] = route
		return
	}

	subMatchs := mux.regexpParamFinded.FindAllStringSubmatch(route.Path, -1)
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
