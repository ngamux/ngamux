package ngamux

import (
	"net/http"
)

type Router interface {
	getPath() string
	HandleFunc(string, string, http.HandlerFunc)
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
	Put(string, http.HandlerFunc)
	Patch(string, http.HandlerFunc)
	Delete(string, http.HandlerFunc)
}

type Annotation struct {
	Mux        Router
	Annotators []Annotator
}
type Annotator func(Route)

func (mux *Ngamux) Annotate(annotators ...Annotator) *Annotation {
	return &Annotation{mux, annotators}
}
func (mux *HttpServeMux) Annotate(annotators ...Annotator) *Annotation {
	return &Annotation{mux, annotators}
}

func (a *Annotation) annotate(method, path string, handler http.HandlerFunc) {
	for _, ann := range a.Annotators {
		if parent := a.Mux; parent != nil {
			ann(buildRoute(parent.getPath()+path, method, handler))
			continue
		}
		ann(buildRoute(path, method, handler))
	}
}

func (a *Annotation) HandlerFunc(method, path string, handler http.HandlerFunc) {
	a.annotate(method, path, handler)
	a.Mux.HandleFunc(method, path, handler)
}
func (a *Annotation) Get(path string, handler http.HandlerFunc) {
	a.annotate(http.MethodGet, path, handler)
	a.Mux.Get(path, handler)
}
func (a *Annotation) Post(path string, handler http.HandlerFunc) {
	a.annotate(http.MethodPost, path, handler)
	a.Mux.Post(path, handler)
}
func (a *Annotation) Put(path string, handler http.HandlerFunc) {
	a.annotate(http.MethodPut, path, handler)
	a.Mux.Put(path, handler)
}
func (a *Annotation) Patch(path string, handler http.HandlerFunc) {
	a.annotate(http.MethodPatch, path, handler)
	a.Mux.Patch(path, handler)
}
func (a *Annotation) Delete(path string, handler http.HandlerFunc) {
	a.annotate(http.MethodDelete, path, handler)
	a.Mux.Delete(path, handler)
}
