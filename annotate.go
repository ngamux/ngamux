package ngamux

import "net/http"

type Router interface {
	HandlerFunc(string, string, Handler)
	Get(string, Handler)
	Post(string, Handler)
	Put(string, Handler)
	Patch(string, Handler)
	Delete(string, Handler)
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

func (a *Annotation) annotate(method, path string, handler Handler) {
	for _, ann := range a.Annotators {
		ann(buildRoute(path, method, handler))
	}
}

func (a *Annotation) HandlerFunc(method, path string, handler Handler) {
	a.annotate(method, path, handler)
	a.Mux.HandlerFunc(method, path, handler)
}
func (a *Annotation) Get(path string, handler Handler) {
	a.annotate(http.MethodGet, path, handler)
	a.Mux.Get(path, handler)
}
func (a *Annotation) Post(path string, handler Handler) {
	a.annotate(http.MethodPost, path, handler)
	a.Mux.Post(path, handler)
}
func (a *Annotation) Put(path string, handler Handler) {
	a.annotate(http.MethodPut, path, handler)
	a.Mux.Put(path, handler)
}
func (a *Annotation) Patch(path string, handler Handler) {
	a.annotate(http.MethodPatch, path, handler)
	a.Mux.Patch(path, handler)
}
func (a *Annotation) Delete(path string, handler Handler) {
	a.annotate(http.MethodDelete, path, handler)
	a.Mux.Delete(path, handler)
}
