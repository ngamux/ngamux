package ngamux

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type (
	Route struct {
		Path       string
		Method     string
		Handler    Handler
		Params     [][]string
		UrlMathcer *regexp.Regexp
	}

	routeMap map[string]Route
)

func buildRoute(url string, status string, handler Handler, middlewares ...MiddlewareFunc) Route {
	handler = WithMiddlewares(middlewares...)(handler)

	return Route{
		Path:    url,
		Handler: handler,
	}
}

func (mux *Ngamux) addRoute(route Route) {
	var (
		err            error
		pathWithParams string
	)

	// check if route doesn't have url param
	if !strings.Contains(route.Path, ":") {
		mux.routes[route.Path] = route
		return
	}

	subMatchs := mux.regexpParamFinded.FindAllStringSubmatch(route.Path, -1)
	route.Params = [][]string{}
	for _, val := range subMatchs {
		route.Params = append(route.Params, []string{val[0][1:]})
	}

	pathWithParams = mux.regexpParamFinded.ReplaceAllString(route.Path, "([0-9a-zA-Z\\.\\-_]+)")
	route.Path = pathWithParams

	route.UrlMathcer, err = regexp.Compile("^" + pathWithParams + "$")
	if err != nil {
		log.Fatal(err)
		return
	}

	mux.routesParam[route.Path] = route
}

func (mux *Ngamux) getRoute(r *http.Request) (Route, *http.Request) {
	path := r.URL.Path
	if mux.config.RemoveTrailingSlash && path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimRight(path, "/")
	}

	foundRoute, ok := mux.routes[path]
	if !ok {
		for url, route := range mux.routesParam {

			if route.UrlMathcer.MatchString(path) {
				foundParams := route.UrlMathcer.FindAllStringSubmatch(path, -1)
				params := make([][]string, len(route.Params))
				copy(params, route.Params)
				for i := range params {
					params[i] = append(params[i], foundParams[0][i+1])
				}
				if len(route.Params) > 0 {
					route.Params = params
					ctx := context.WithValue(r.Context(), KeyContextParams, params)
					r = r.WithContext(ctx)
				}
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

	return foundRoute, r
}
