package ngamux

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type KeyContext int

const (
	KeyContextParams KeyContext = 1 << iota
)

type Config struct {
	RemoveTrailingSlash bool
	NotFoundHandler     http.HandlerFunc
	PanicHandler        func(w http.ResponseWriter, r *http.Request, p interface{})
}

type Ngamux struct {
	routes            map[string]map[string]Route
	routesParam       map[string]map[string]Route
	config            Config
	regexpParamFinded *regexp.Regexp
}

var _ http.Handler = &Ngamux{}

type Route struct {
	Path       string
	Handler    http.HandlerFunc
	Params     [][]string
	UrlMathcer *regexp.Regexp
}

func handlerNotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(rw, "404 page not found")
}

func DefaultPanicHandler(rw http.ResponseWriter, _ *http.Request, p interface{}) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, "panic: %v", p)
}

func buildRoute(url string, handler http.HandlerFunc) Route {
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

	if config.PanicHandler == nil {
		config.PanicHandler = DefaultPanicHandler
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

func (r *Ngamux) addRoute(method string, route Route) {
	var (
		err            error
		pathWithParams string
		subMatchs      = r.regexpParamFinded.FindAllStringSubmatch(route.Path, -1)
	)

	route.Params = [][]string{}

	// check if route not have url param
	if len(subMatchs) == 0 {
		r.routes[method][route.Path] = route
		return
	}

	for _, val := range subMatchs {
		route.Params = append(route.Params, []string{val[0][1:]})
	}

	pathWithParams = r.regexpParamFinded.ReplaceAllString(route.Path, "([0-9a-zA-Z]+)")
	route.Path = pathWithParams

	route.UrlMathcer, err = regexp.Compile("^" + pathWithParams + "$")
	if err != nil {
		log.Fatal(err)
		return
	}

	r.routesParam[method][route.Path] = route
}

func (r *Ngamux) getRoute(method string, path string) Route {
	foundRoute := Route{
		Handler: r.config.NotFoundHandler,
	}

	foundRoute, ok := r.routes[method][path]
	if !ok {
		for url, route := range r.routesParam[method] {

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

	return foundRoute
}

func (mux *Ngamux) Group(path string, middlewares ...http.HandlerFunc) *group {
	group := &group{
		parent:      mux,
		middlewares: middlewares,
		path:        path,
	}
	return group
}

func (mux *Ngamux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			mux.config.PanicHandler(rw, r, p)
		}
	}()

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

func (mux *Ngamux) Get(path string, handler http.HandlerFunc) {
	mux.addRoute(http.MethodGet, buildRoute(path, handler))
}

func (mux *Ngamux) Post(path string, handler http.HandlerFunc) {
	mux.addRoute(http.MethodPost, buildRoute(path, handler))
}

func (mux *Ngamux) Patch(path string, handler http.HandlerFunc) {
	mux.addRoute(http.MethodPatch, buildRoute(path, handler))
}

func (mux *Ngamux) Put(path string, handler http.HandlerFunc) {
	mux.addRoute(http.MethodPut, buildRoute(path, handler))
}

func (mux *Ngamux) Delete(path string, handler http.HandlerFunc) {
	mux.addRoute(http.MethodDelete, buildRoute(path, handler))
}
