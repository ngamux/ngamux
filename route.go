package ngamux

import (
	"fmt"
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
		key      string
		path     []byte
		param    string
		handler  http.Handler
		children map[string]*Node
	}
)

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
		current = &Node{key: "", children: make(map[string]*Node)}
		t.root.Set(method, current)
	}

	keys := strings.Split(key, "/")
	path := []byte{}
	for i, k := range keys {
		if i > 0 {
			path = append(path, '/')
			path = append(path, []byte(k)...)
		}
		isWildcard := strings.HasPrefix(k, "{") && strings.HasSuffix(k, "}")
		var param string
		if isWildcard {
			param = k[1 : len(k)-1]
			k = "{}"
		}
		if _, ok := current.children[k]; !ok {
			current.children[k] = &Node{key: k, children: make(map[string]*Node)}
			current.children[k].param = param
		}
		current = current.children[k]
	}
	current.handler = handler
	current.path = path
}

func (t Ngamux) match(key string, params map[string]string, handler *http.Handler, pattern *string) {
	method, key := splitMethodPath(key)
	current, ok := t.root.Get(method)
	if !ok {
		t.root.Each(func(k string, n *Node) bool {
			matchNode(n, key, params, handler, pattern)
			if handler != nil && k != "ALL" {
				*handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = fmt.Fprintf(w, "%d method not allowed\n", http.StatusMethodNotAllowed)
				})
				return false
			}
			return true
		})
		return
	}

	matchNode(current, key, params, handler, pattern)
}

func matchNode(current *Node, key string, params map[string]string, handler *http.Handler, pattern *string) {
	keys := strings.Split(key, "/")
	for _, k := range keys {
		if _, ok := current.children[k]; ok {
			current = current.children[k]
		} else if _, ok := current.children["{}"]; ok {
			current = current.children["{}"]
			params[current.param] = k
		} else {
			return
		}
	}

	*handler = current.handler
	*pattern = string(current.path)
}

func (t Ngamux) Handler(r *http.Request) (http.Handler, string) {
	params := make(map[string]string)
	var handler http.Handler
	var pattern string
	t.match(r.Method+" "+r.URL.Path, params, &handler, &pattern)
	if handler == nil {
		return nil, ""
	}

	for k, v := range params {
		r.SetPathValue(k, v)
	}
	return handler, pattern
}

func (t Ngamux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, _ := t.Handler(r)
	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	r.Method = "ALL"
	handler, _ = t.Handler(r)
	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}
