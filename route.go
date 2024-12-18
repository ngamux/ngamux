package ngamux

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type (
	// Route describe a route object
	Route struct {
		RawPath    string
		Path       string
		Method     string
		Handler    http.HandlerFunc
		Params     [][]string
		URLMatcher *regexp.Regexp
	}

	routeMap map[string]map[string]Route
)

func buildRoute(url string, method string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) Route {
	handler = WithMiddlewares(middlewares...)(handler)

	route := Route{
		RawPath: url,
		Path:    url,
		Method:  method,
		Handler: handler,
	}

	return route
}

func (mux *Ngamux) addRoute(route Route) {
	var (
		err            error
		pathWithParams string
	)

	// check if route doesn't have url param
	if !strings.Contains(route.Path, ":") && !strings.Contains(route.Path, "+") {
		if mux.routes[route.Path] == nil {
			mux.routes[route.Path] = map[string]Route{}
		}

		mux.routes[route.Path][route.Method] = route
		mux.Log(LogLevelInfo, "[ROUTE] %s %s registered", route.Method, route.RawPath)
		return
	}

	// building route with url param
	subMatchs := mux.regexpParamFinded.FindAllStringSubmatch(route.Path, -1)
	route.Params = [][]string{}
	for _, val := range subMatchs {
		route.Params = append(route.Params, []string{val[0][1:]})
	}

	if strings.Contains(route.Path, ":") {
		pathWithParams = mux.regexpParamFinded.ReplaceAllString(route.Path, "([0-9a-zA-Z\\.\\-_]+)")
	} else if strings.Contains(route.Path, "+") {
		pathWithParams = mux.regexpParamFinded.ReplaceAllString(route.Path, "([0-9a-zA-Z\\.\\-/_]+)")
	}
	route.Path = pathWithParams

	route.URLMatcher, err = regexp.Compile("^" + pathWithParams + "$")
	if err != nil {
		log.Fatal(err)
		return
	}

	if mux.routesParam[route.Path] == nil {
		mux.routesParam[route.Path] = map[string]Route{}
	}
	mux.routesParam[route.Path][route.Method] = route
	mux.Log(LogLevelInfo, "[ROUTE] %s %s registered", route.Method, route.RawPath)
}

func buildURLParams(r *http.Request, route Route, path string) (*http.Request, Route) {
	if route.URLMatcher == nil {
		return r, route
	}

	foundParams := route.URLMatcher.FindAllStringSubmatch(path, -1)
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
			URLMatcher, err := regexp.Compile("^" + url + "$")
			if err != nil {
				break
			}

			if URLMatcher.MatchString(path) {
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
		tmpR := Req(r)
		tmpR.Locals("error", ErrorNotFound)
		r = tmpR.Request
		foundRoute.Handler = mux.config.GlobalErrorHandler
	} else {
		route, ok := foundRouteMap[r.Method]
		if !ok {
			tmpR := Req(r)
			tmpR.Locals("error", ErrorMethodNotAllowed)
			r = tmpR.Request
			foundRoute.Handler = mux.config.GlobalErrorHandler
		} else {

			if r.Method == http.MethodHead {
				r.Body = nil
			}

			r, route = buildURLParams(r, route, path)
			foundRoute = route
		}
	}

	return foundRoute, r
}
