package ngamux

// Group returns new nested ngamux object
func (mux *Ngamux) Group(url string) *Ngamux {
	group := New()
	group.path = url
	group.parent = mux
	return group
}

// func (mux *Ngamux) addRouteFromGroup(route Route) {
// 	url := path.Join(mux.path, route.Path)
// 	middlewares := mux.middlewares
// 	parent := mux.parent
// 	for parent != nil {
// 		url = path.Join(parent.path, url)
// 		middlewares = append(middlewares, parent.middlewares...)
// 		if parent.parent == nil {
// 			break
// 		}
//
// 		parent = parent.parent
// 	}
//
// 	parent.addRoute(buildRoute(url, route.Method, route.Handler, middlewares...))
// }
