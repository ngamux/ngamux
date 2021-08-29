package ngamux

import (
	"net/http"
	"regexp"
)

type Router struct {
	routes map[string]map[string]Route
	config Config
}

type Route struct {
	Path     string
	Handlers []http.HandlerFunc
	Params   [][]string
}

func newRouter(config Config) *Router {
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

	return &Router{
		routes: routesMap,
		config: config,
	}
}

func (r *Router) AddRoute(method string, route Route) {
	paramsFinder, err := regexp.Compile("(:[a-zA-Z][0-9a-zA-Z]+)")
	if err != nil {
		r.routes[method][route.Path] = route
	}

	route.Params = [][]string{}
	for _, val := range paramsFinder.FindAllStringSubmatch(route.Path, -1) {
		route.Params = append(route.Params, []string{val[0][1:]})
	}

	pathWithParams := paramsFinder.ReplaceAllString(route.Path, "([0-9a-zA-Z]+)")
	route.Path = pathWithParams

	r.routes[method][route.Path] = route
}

func (r *Router) GetRoute(method string, path string) Route {
	foundRoute := Route{
		Handlers: []http.HandlerFunc{r.config.NotFoundHandler},
	}

	foundRoute, ok := r.routes[method][path]
	if !ok {
		for url, route := range r.routes[method] {
			urlMatcher, err := regexp.Compile("^" + url + "$")
			if err != nil {
				continue
			}

			if urlMatcher.MatchString(path) {
				foundParams := urlMatcher.FindAllStringSubmatch(path, -1)
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
