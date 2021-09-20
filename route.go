package ngamux

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

type (
	Route struct {
		Path       string
		Handler    HandlerFunc
		Params     [][]string
		UrlMathcer *regexp.Regexp
	}

	routeMap map[string]map[string]Route
)

func buildRouteMap() routeMap {
	return routeMap{
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
}

func buildRoute(url string, handler HandlerFunc, middlewares ...MiddlewareFunc) Route {
	handler = WithMiddlewares(middlewares...)(handler)

	return Route{
		Path:    url,
		Handler: handler,
	}
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
