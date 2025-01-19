package ngamux

import (
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

func (t *Ngamux) handle(key string, handler http.Handler) {
	current := t.root
	keys := strings.Split(key, "/")
	for _, k := range keys {
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
}

func (t Ngamux) match(key string, params map[string]string) (http.Handler, string) {
	pattern := strings.Builder{}
	current := t.root
	keys := strings.Split(key, "/")
	for _, k := range keys {
		if _, ok := current.children[k]; ok {
			current = current.children[k]
			pattern.WriteString("/")
			pattern.WriteString(k)
		} else if _, ok := current.children["{}"]; ok {
			current = current.children["{}"]
			params[current.param] = k

			pattern.WriteString("/")
			pattern.WriteString("{")
			pattern.WriteString(current.param)
			pattern.WriteString("}")
		} else {
			return nil, ""
		}
	}

	_, pattern1, _ := strings.Cut(pattern.String(), " ")
	return current.handler, pattern1
}

func (t Ngamux) Handler(r *http.Request) (http.Handler, string) {
	params := make(map[string]string)
	handler, pattern := t.match(r.Method+" "+r.URL.Path, params)
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
