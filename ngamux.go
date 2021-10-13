package ngamux

import (
	"errors"
	"fmt"
	"github.com/ngamux/ngamux/constants"
	ngamuxerr "github.com/ngamux/ngamux/errors"
	"net/http"
	"regexp"
)

type KeyContext int

const (
	KeyContextParams KeyContext = 1 << iota
)

type (
	// Ngamux main components struct defined inside Ngamux
	Ngamux struct {
		parent             *Ngamux
		path               string
		routes             routeMap
		routesParam        routeMap
		config             Config
		regexpParamFounded *regexp.Regexp
		middlewares        []MiddlewareFunc
	}

	// INgamux representations of *Ngamux functionality
	INgamux interface {
		// ServeHTTP overrides from http package
		ServeHTTP(rw http.ResponseWriter, r *http.Request)

		// Use attaching or appending middleware to your route endpoint as many as you want
		Use(middlewares ...MiddlewareFunc)

		// Get create route endpoint with http get method operation
		// It returns nothing
		Get(url string, handler Handler)

		// Post create route endpoint with http post method operation
		// It returns nothing
		Post(url string, handler Handler)

		// Patch create route endpoint with http patch method operation
		// It returns nothing
		Patch(url string, handler Handler)

		// Put create route endpoint with http put method operation
		// It returns nothing
		Put(url string, handler Handler)

		// Delete create route endpoint with http delete method operation
		// It returns nothing
		Delete(url string, handler Handler)

		// All create route endpoint with all http methods operation
		// It returns nothing
		All(url string, handler Handler)

		// Group create group of endpoints with specific customizable middleware
		// It returns *Ngamux
		Group(url string, middlewares ...MiddlewareFunc) *Ngamux

		ngamux() *Ngamux
		addRoute(route Route)
		getRoute(r *http.Request) (Route, *http.Request)
	}
)

var (
	_ http.Handler = &Ngamux{}
	_ http.Handler = Handler(func(rw http.ResponseWriter, r *http.Request) error { return nil })

	paramsFinder       = regexp.MustCompile("(:[a-zA-Z]+[0-9a-zA-Z]*)")
	globalErrorHandler = func(rw http.ResponseWriter, r *http.Request) error {
		err := GetContextValue(r, constants.ContextKeyError).(error)

		if errors.Is(err, ngamuxerr.NotFound) {
			rw.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, ngamuxerr.MethodNotAllowed) {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}

		_, _ = fmt.Fprintln(rw, err)

		return nil
	}
)

func NewNgamux(configs ...Config) INgamux {
	config := buildConfig(configs...)
	routesMap := routeMap{}
	routesParamMap := routeMap{}
	router := &Ngamux{
		routes:             routesMap,
		routesParam:        routesParamMap,
		config:             config,
		regexpParamFounded: paramsFinder,
	}

	return router
}

func (mux *Ngamux) ngamux() *Ngamux {
	return mux
}

func (mux *Ngamux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, r := mux.getRoute(r)
	err := route.Handler(rw, r)
	if err != nil {
		rw.WriteHeader(500)
		_, _ = rw.Write([]byte(err.Error()))
	}
}

func (mux *Ngamux) Use(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
	mux.config.GlobalErrorHandler = WithMiddlewares(mux.middlewares...)(mux.config.GlobalErrorHandler)
}

func (mux *Ngamux) Get(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodGet, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodGet, handler, mux.middlewares...))
}

func (mux *Ngamux) Post(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPost, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPost, handler, mux.middlewares...))
}

func (mux *Ngamux) Patch(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPatch, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPatch, handler, mux.middlewares...))
}

func (mux *Ngamux) Put(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodPut, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodPut, handler, mux.middlewares...))
}

func (mux *Ngamux) Delete(url string, handler Handler) {
	if mux.parent != nil {
		mux.addRouteFromGroup(buildRoute(url, http.MethodDelete, handler))
		return
	}
	mux.addRoute(buildRoute(url, http.MethodDelete, handler, mux.middlewares...))
}

func (mux *Ngamux) All(url string, handler Handler) {
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		if mux.parent != nil {
			mux.addRouteFromGroup(buildRoute(url, method, handler))
			return
		}

		mux.addRoute(buildRoute(url, method, handler, mux.middlewares...))
	}
}
