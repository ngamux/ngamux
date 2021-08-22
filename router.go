package ngamux

import (
	"fmt"
	"net/http"

	"github.com/ngamux/gotrie"
)

type Router struct {
	routes map[string]*gotrie.Trie
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

func newRouter() *Router {
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
	}
}

func (r *Router) AddRoute(method string, route Route) {
	r.routes[method].Put(route.Path, route)
}

func (r *Router) GetRoute(method string, path string) Route {
	route := r.routes[method].Get(path)
	if route == nil {
		return Route{
			Handler: func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(rw, "404 page not found")
			},
		}
	}

	routeFound := r.routes[method].Get(path).(Route)
	return routeFound
}
