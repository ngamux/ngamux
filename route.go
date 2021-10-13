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
		UrlMatcher *regexp.Regexp
	}

	routeMap map[string]map[string]Route
)

func buildRoute(url string, method string, handler Handler, middlewares ...MiddlewareFunc) Route {
	handler = WithMiddlewares(middlewares...)(handler)

	return Route{
		Path:    url,
		Method:  method,
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
		if mux.routes[route.Path] == nil {
			mux.routes[route.Path] = map[string]Route{}
		}

		mux.routes[route.Path][route.Method] = route
		return
	}

	// building route with url param
	subMatchs := mux.regexpParamFounded.FindAllStringSubmatch(route.Path, -1)
	route.Params = [][]string{}
	for _, val := range subMatchs {
		route.Params = append(route.Params, []string{val[0][1:]})
	}

	pathWithParams = mux.regexpParamFounded.ReplaceAllString(route.Path, "([0-9a-zA-Z\\.\\-_]+)")
	route.Path = pathWithParams

	route.UrlMatcher, err = regexp.Compile("^" + pathWithParams + "$")
	if err != nil {
		log.Fatal(err)
		return
	}

	if mux.routesParam[route.Path] == nil {
		mux.routesParam[route.Path] = map[string]Route{}
	}
	mux.routesParam[route.Path][route.Method] = route
}

func buildUrlParams(r *http.Request, route Route, path string) (*http.Request, Route) {
	if route.UrlMatcher == nil {
		return r, route
	}

	foundParams := route.UrlMatcher.FindAllStringSubmatch(path, -1)
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
	return r, route
}

func (mux *Ngamux) getRoute(r *http.Request) (Route, *http.Request) {
	path := r.URL.Path
	if mux.config.RemoveTrailingSlash && path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimRight(path, "/")
	}

	foundRouteMap, ok := mux.routes[path]
	if !ok {
		for url, route := range mux.routesParam {
			urlMatcher, err := regexp.Compile("^" + url + "$")
			if err != nil {
				break
			}

			if urlMatcher.MatchString(path) {
				foundRouteMap = route
				break
			}

			if url == path {
				foundRouteMap = route
				break
			}
		}
	}

	var foundRoute Route
	if len(foundRouteMap) <= 0 {
		r = SetContextValue(r, "error", ErrorNotFound)
		foundRoute.Handler = mux.config.GlobalErrorHandler
	} else {
		route, ok := foundRouteMap[r.Method]
		if !ok {
			r = SetContextValue(r, "error", ErrorMethodNotAllowed)
			foundRoute.Handler = mux.config.GlobalErrorHandler
		} else {
			r, route = buildUrlParams(r, route, path)
			foundRoute = route
		}
	}

	return foundRoute, r
}
