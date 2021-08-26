package ngamux

import (
	"fmt"
	"net/http"
)

type Ngamux struct {
	config Config
	router *Router
	parent *Route
}

type Config struct {
	NotFoundHandler http.HandlerFunc
}

var _ http.Handler = &Ngamux{}

func NewNgamux(config ...Config) *Ngamux {
	ngamuxConfig := Config{}
	if len(config) > 0 {
		ngamuxConfig = config[0]
	}

	if ngamuxConfig.NotFoundHandler == nil {
		ngamuxConfig.NotFoundHandler = func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(rw, "404 page not found")
		}
	}
	return &Ngamux{
		config: ngamuxConfig,
		router: newRouter(ngamuxConfig),
	}
}

func (mux *Ngamux) Group(path string, middlewares ...http.HandlerFunc) *Ngamux {
	ngamux := NewNgamux(mux.config)
	ngamux.parent = &Route{
		Path:     path,
		Handlers: middlewares,
	}
	return ngamux
}

func (mux *Ngamux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := mux.router.GetRoute(r.Method, r.URL.Path)
	for _, handler := range route.Handlers {
		handler(w, r)
	}
}

func (mux *Ngamux) Get(path string, handler ...http.HandlerFunc) {
	if mux.parent != nil {
		path = fmt.Sprintf("%s%s", mux.parent.Path, path)
	}

	mux.router.AddRoute(http.MethodGet, Route{
		Path:     path,
		Handlers: handler,
	})
}

func (mux *Ngamux) Post(path string, handler ...http.HandlerFunc) {
	if mux.parent != nil {
		path = fmt.Sprintf("%s%s", mux.parent.Path, path)
	}

	mux.router.AddRoute(http.MethodPost, Route{
		Path:     path,
		Handlers: handler,
	})
}

func (mux *Ngamux) Patch(path string, handler ...http.HandlerFunc) {
	if mux.parent != nil {
		path = fmt.Sprintf("%s%s", mux.parent.Path, path)
	}

	mux.router.AddRoute(http.MethodPatch, Route{
		Path:     path,
		Handlers: handler,
	})
}

func (mux *Ngamux) Put(path string, handler ...http.HandlerFunc) {
	if mux.parent != nil {
		path = fmt.Sprintf("%s%s", mux.parent.Path, path)
	}

	mux.router.AddRoute(http.MethodPut, Route{
		Path:     path,
		Handlers: handler,
	})
}

func (mux *Ngamux) Delete(path string, handler ...http.HandlerFunc) {
	if mux.parent != nil {
		path = fmt.Sprintf("%s%s", mux.parent.Path, path)
	}

	mux.router.AddRoute(http.MethodDelete, Route{
		Path:     path,
		Handlers: handler,
	})
}
