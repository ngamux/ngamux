package ngamux

import (
	"fmt"
	"net/http"
)

type Ngamux struct {
	router *Router
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
		router: newRouter(ngamuxConfig),
	}
}

func (mux *Ngamux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := mux.router.GetRoute(r.Method, r.URL.Path)
	route.Handler(w, r)
}

func (mux *Ngamux) Get(path string, handler http.HandlerFunc) {
	mux.router.AddRoute(http.MethodGet, Route{
		Path:    path,
		Handler: handler,
	})
}

func (mux *Ngamux) Post(path string, handler http.HandlerFunc) {
	mux.router.AddRoute(http.MethodPost, Route{
		Path:    path,
		Handler: handler,
	})
}

func (mux *Ngamux) Patch(path string, handler http.HandlerFunc) {
	mux.router.AddRoute(http.MethodPatch, Route{
		Path:    path,
		Handler: handler,
	})
}

func (mux *Ngamux) Put(path string, handler http.HandlerFunc) {
	mux.router.AddRoute(http.MethodPut, Route{
		Path:    path,
		Handler: handler,
	})
}

func (mux *Ngamux) Delete(path string, handler http.HandlerFunc) {
	mux.router.AddRoute(http.MethodDelete, Route{
		Path:    path,
		Handler: handler,
	})
}
