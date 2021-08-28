package ngamux

import (
	"net/http"
)

type Router struct {
	routes map[string]map[string]Route
	config Config
}

type Route struct {
	Path     string
	Handlers []http.HandlerFunc
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
	r.routes[method][route.Path] = route
}

func (r *Router) GetRoute(method string, path string) Route {
	route, ok := r.routes[method][path]
	if !ok {
		return Route{
			Handlers: []http.HandlerFunc{r.config.NotFoundHandler},
		}
	}

	return route
}
