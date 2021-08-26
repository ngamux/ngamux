package ngamux

import (
	"net/http"

	"github.com/ngamux/gotrie"
)

type Router struct {
	routes map[string]*gotrie.Trie
	config Config
}

type Route struct {
	Path     string
	Handlers []http.HandlerFunc
}

func newRouter(config Config) *Router {
	routesMap := map[string]*gotrie.Trie{
		http.MethodGet:     gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodPost:    gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodPatch:   gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodPut:     gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodDelete:  gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodOptions: gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodConnect: gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodHead:    gotrie.NewTrie(gotrie.Config{Separator: "/"}),
		http.MethodTrace:   gotrie.NewTrie(gotrie.Config{Separator: "/"}),
	}

	return &Router{
		routes: routesMap,
		config: config,
	}
}

func (r *Router) AddRoute(method string, route Route) {
	r.routes[method].Put(route.Path, route)
}

func (r *Router) GetRoute(method string, path string) Route {
	route := r.routes[method].Get(path)
	if route == nil {
		return Route{
			Handlers: []http.HandlerFunc{r.config.NotFoundHandler},
		}
	}

	routeFound := r.routes[method].Get(path).(Route)
	return routeFound
}
