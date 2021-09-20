package ngamux

import "path"

func (mux *Ngamux) Group(url string, middlewares ...MiddlewareFunc) *Ngamux {
	group := &Ngamux{
		parent:      mux,
		path:        url,
		middlewares: middlewares,
	}
	return group
}

func (mux *Ngamux) addRouteFromGroup(method string, route Route) {
	url := path.Join(mux.path, route.Path)
	middlewares := mux.middlewares
	middlewares = append(middlewares, mux.parent.middlewares...)
	mux.parent.addRoute(method, buildRoute(url, route.Handler, middlewares...))
}
