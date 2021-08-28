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
		http.MethodGet:     map[string]Route{},
		http.MethodPost:    map[string]Route{},
		http.MethodPatch:   map[string]Route{},
		http.MethodPut:     map[string]Route{},
		http.MethodDelete:  map[string]Route{},
		http.MethodOptions: map[string]Route{},
		http.MethodConnect: map[string]Route{},
		http.MethodHead:    map[string]Route{},
		http.MethodTrace:   map[string]Route{},
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
