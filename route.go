package ngamux

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

type (
	Route struct {
		RawPath    string
		Path       string
		Method     string
		Handler    http.HandlerFunc
		Params     [][]string
		URLMatcher *regexp.Regexp
	}

	Node struct {
		path     string
		param    string
		handler  http.Handler
		children map[string]*Node
	}
)

const maxStackRouteParams = 8

type routeParam struct {
	name  string
	value string
}

type routeParams struct {
	small [maxStackRouteParams]routeParam
	count int
	extra map[string]string
}

func (p *routeParams) reset() {
	p.count = 0
	if p.extra != nil {
		clear(p.extra)
	}
}

func (p *routeParams) add(name, value string) {
	if p.count < len(p.small) {
		p.small[p.count] = routeParam{name: name, value: value}
		p.count++
		return
	}

	if p.extra == nil {
		p.extra = make(map[string]string)
	}
	p.extra[name] = value
}

func (p *routeParams) set(r *http.Request) {
	for i := 0; i < p.count; i++ {
		r.SetPathValue(p.small[i].name, p.small[i].value)
	}
	for name, value := range p.extra {
		r.SetPathValue(name, value)
	}
}

var methodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = io.WriteString(w, "405 method not allowed\n")
})

func (t *Ngamux) Handle(key string, handler http.Handler) {
	if strings.HasPrefix(key, "/") {
		key = "ALL " + key
	}
	t.handle(key, handler)
}

func splitMethodPath(path string) (string, string) {
	paths := strings.Split(path, " ")
	method := "ALL"
	if len(path) > 1 {
		method = paths[0]
		path = paths[1]
	}
	return method, path
}

func (t *Ngamux) handle(key string, handler http.Handler) {
	method, key := splitMethodPath(key)
	current, ok := t.root.Get(method)
	if !ok {
		current = &Node{children: make(map[string]*Node)}
		t.root.Set(method, current)
	}

	keys := strings.Split(key, "/")
	var pattern strings.Builder
	pattern.Grow(len(key))
	for i, k := range keys {
		if i > 0 {
			pattern.WriteByte('/')
			pattern.WriteString(k)
		}
		isWildcard := strings.HasPrefix(k, "{") && strings.HasSuffix(k, "}")
		var param string
		if isWildcard {
			param = k[1 : len(k)-1]
			k = "{}"
		}
		if _, ok := current.children[k]; !ok {
			current.children[k] = &Node{children: make(map[string]*Node)}
			current.children[k].param = param
		}
		current = current.children[k]
	}
	current.handler = handler
	current.path = pattern.String()
}

func (t *Ngamux) matchMethod(method, path string, params *routeParams) (http.Handler, string, bool) {
	current, ok := t.root.Get(method)
	if !ok {
		return nil, "", false
	}

	node, matched := matchNode(current, path, params)
	if !matched || node.handler == nil {
		return nil, "", true
	}
	return node.handler, node.path, true
}

func matchNode(current *Node, path string, params *routeParams) (*Node, bool) {
	for {
		segment := path
		separator := strings.IndexByte(path, '/')
		if separator >= 0 {
			segment = path[:separator]
			path = path[separator+1:]
		}

		next, ok := current.children[segment]
		if !ok {
			next, ok = current.children["{}"]
			if !ok {
				return nil, false
			}
			params.add(next.param, segment)
		}
		current = next

		if separator < 0 {
			return current, true
		}
	}
}

func (t *Ngamux) matchOtherMethod(method, path string, params *routeParams) (http.Handler, string) {
	var handler http.Handler
	var pattern string
	t.root.Each(func(k string, n *Node) bool {
		if k == method || k == "ALL" {
			return true
		}

		params.reset()
		node, matched := matchNode(n, path, params)
		if matched && node.handler != nil {
			handler = methodNotAllowedHandler
			pattern = node.path
			return false
		}
		return true
	})
	return handler, pattern
}

func (t Ngamux) Handler(r *http.Request) (http.Handler, string) {
	var params routeParams
	handler, pattern, rootFound := t.matchMethod(r.Method, r.URL.Path, &params)
	if handler != nil {
		params.set(r)
		return handler, pattern
	}
	if rootFound {
		return nil, ""
	}

	params.reset()
	if r.Method != "ALL" {
		handler, pattern, _ = t.matchMethod("ALL", r.URL.Path, &params)
		if handler != nil {
			params.set(r)
			return handler, pattern
		}
	}

	handler, pattern = t.matchOtherMethod(r.Method, r.URL.Path, &params)
	if handler != nil {
		params.set(r)
	}
	return handler, pattern
}

func (t Ngamux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params routeParams
	handler, _, _ := t.matchMethod(r.Method, r.URL.Path, &params)
	if handler == nil && r.Method != "ALL" {
		params.reset()
		handler, _, _ = t.matchMethod("ALL", r.URL.Path, &params)
	}

	if handler == nil {
		if _, ok := t.root.Get("ALL"); !ok {
			handler, _ = t.matchOtherMethod(r.Method, r.URL.Path, &params)
		}
	}

	if handler == nil {
		http.NotFound(w, r)
		return
	}

	params.set(r)
	handler.ServeHTTP(w, r)
}
