package ngamux

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ngamux/ngamux/middleware"
)

type KeyContext int

const (
	KeyContextParams KeyContext = 1 << iota
)

type Ngamux struct {
	config   Config
	router   *Router
	recovery *middleware.Recovery
	cors     *middleware.Cors
}

type Config struct {
	RemoveTrailingSlash bool
	NotFoundHandler     http.HandlerFunc
}

var _ http.Handler = &Ngamux{}

func handlerNotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(rw, "404 page not found")
}

func buildRoute(url string, handler ...http.HandlerFunc) Route {
	return Route{
		Path:     url,
		Handlers: handler,
	}
}

func NewNgamux(config ...Config) *Ngamux {
	ngamuxConfig := Config{
		RemoveTrailingSlash: true,
	}
	if len(config) > 0 {
		ngamuxConfig = config[0]
	}

	if ngamuxConfig.NotFoundHandler == nil {
		ngamuxConfig.NotFoundHandler = handlerNotFound
	}

	return &Ngamux{
		config: ngamuxConfig,
		router: newRouter(ngamuxConfig),
	}
}

func (mux *Ngamux) Group(path string, middlewares ...http.HandlerFunc) *group {
	group := &group{
		parent:      mux,
		middlewares: middlewares,
		path:        path,
	}
	return group
}

func (mux *Ngamux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if mux.config.RemoveTrailingSlash && len(url) > 1 && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	route := mux.router.GetRoute(r.Method, url)

	if len(route.Params) > 0 {
		ctx := context.WithValue(r.Context(), KeyContextParams, route.Params)
		r = r.WithContext(ctx)
	}

	for _, handler := range route.Handlers {
		if mux.recovery != nil {
			defer mux.recovery.RecoverFunc()
		}
		if mux.cors != nil {
			mux.useCors(&w, r)
		}
		handler(w, r)
	}
}

func (mux *Ngamux) useCors(w *http.ResponseWriter, r *http.Request) {

	(*w).Header().Set("Access-Control-Allow-Origin", mux.cors.Origin)
	(*w).Header().Set("Access-Control-Allow-Methods", mux.cors.Method)
	(*w).Header().Set("Access-Control-Allow-Headers", mux.cors.Header)
	(*w).Header().Set("Access-Control-Expose-Headers", mux.cors.ExposeHeader)

}

func (mux *Ngamux) Use(middlewares ...interface{}) {

	for _, midd := range middlewares {

		recovery, ok := midd.(*middleware.Recovery)
		if ok {
			mux.recovery = recovery
			continue
		}

		cors, ok := midd.(*middleware.Cors)
		if ok {
			mux.cors = cors
			continue
		}
	}
}

func (mux *Ngamux) Get(path string, handler ...http.HandlerFunc) {
	mux.router.AddRoute(http.MethodGet, buildRoute(path, handler...))
}

func (mux *Ngamux) Post(path string, handler ...http.HandlerFunc) {
	mux.router.AddRoute(http.MethodPost, buildRoute(path, handler...))
}

func (mux *Ngamux) Patch(path string, handler ...http.HandlerFunc) {
	mux.router.AddRoute(http.MethodPatch, buildRoute(path, handler...))
}

func (mux *Ngamux) Put(path string, handler ...http.HandlerFunc) {
	mux.router.AddRoute(http.MethodPut, buildRoute(path, handler...))
}

func (mux *Ngamux) Delete(path string, handler ...http.HandlerFunc) {
	mux.router.AddRoute(http.MethodDelete, buildRoute(path, handler...))
}
