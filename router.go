package ngamux

import (
	"log"
	"net/http"
	"regexp"
)

type Router struct {
	routes            map[string]map[string]Route
	routesParam       map[string]map[string]Route
	config            Config
	regexpParamFinded *regexp.Regexp
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Params  [][]string
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

	paramsFinder, err := regexp.Compile("(:[a-zA-Z][0-9a-zA-Z]+)")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	router := &Router{
		routes:            routesMap,
		routesParam:       routesMap,
		config:            config,
		regexpParamFinded: paramsFinder,
	}

	for key, val := range router.routes {
		router.routesParam[key] = val
	}

	return router
}

func (r *Router) AddRoute(method string, route Route) {
	var (
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
	r.routesParam[method][route.Path] = route
}

func (r *Router) GetRoute(method string, path string) Route {
	foundRoute := Route{
		Handler: r.config.NotFoundHandler,
	}

	foundRoute, ok := r.routes[method][path]
	if !ok {
		for url, route := range r.routesParam[method] {
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
