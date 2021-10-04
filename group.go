package ngamux

import (
	"path"
)

func (mux *Ngamux) Group(url string, middlewares ...MiddlewareFunc) *Ngamux {
	if mux.parent != nil {
		panic("nested route group is not supported yet")
	}

	group := &Ngamux{
		parent:      mux,
		path:        url,
		middlewares: middlewares,
	}
	return group
}

func (mux *Ngamux) addRouteFromGroup(route Route) {
	url := path.Join(mux.path, route.Path)
	middlewares := mux.middlewares
	middlewares = append(middlewares, mux.parent.middlewares...)
	mux.parent.addRoute(buildRoute(url, route.Method, route.Handler, middlewares...))
}
