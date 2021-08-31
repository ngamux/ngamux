package ngamux

import (
	"net/http"
	"path"
)

type group struct {
	parent      *Ngamux
	path        string
	middlewares []http.HandlerFunc
}

func (mux *group) Get(url string, handler http.HandlerFunc) {
	url = path.Join(mux.path, url)
	mux.parent.Get(url, handler)
}

func (mux *group) Post(url string, handler http.HandlerFunc) {
	url = path.Join(mux.path, url)
	mux.parent.Post(url, handler)
}

func (mux *group) Patch(url string, handler http.HandlerFunc) {
	url = path.Join(mux.path, url)
	mux.parent.Patch(url, handler)
}

func (mux *group) Put(url string, handler http.HandlerFunc) {
	url = path.Join(mux.path, url)
	mux.parent.Put(url, handler)
}

func (mux *group) Delete(url string, handler http.HandlerFunc) {
	url = path.Join(mux.path, url)
	mux.parent.Delete(url, handler)
}
