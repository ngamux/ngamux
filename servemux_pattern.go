package ngamux

import (
	"net/http"
	"net/url"
	"strings"
)

type serveMuxPatternPartKind uint8

const (
	serveMuxLiteral serveMuxPatternPartKind = iota
	serveMuxWildcard
	serveMuxMultiWildcard
	serveMuxEnd
	serveMuxTrailingSlash
)

type serveMuxPatternPart struct {
	kind serveMuxPatternPartKind
	name string
}

func wrapServeMuxHandler(pattern string, next http.HandlerFunc) http.HandlerFunc {
	parts := parseServeMuxPattern(pattern)
	if len(parts) == 0 {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setServeMuxPathValues(r, parts)
		next(w, r)
	}
}

func registerServeMuxHandler(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	if handler == nil {
		mux.HandleFunc(pattern, nil)
		return
	}
	mux.HandleFunc(pattern, wrapServeMuxHandler(pattern, handler))
}

func parseServeMuxPattern(pattern string) []serveMuxPatternPart {
	path := pattern
	if i := strings.IndexAny(pattern, " \t"); i >= 0 {
		path = strings.TrimLeft(pattern[i+1:], " \t")
	}
	if !strings.HasPrefix(path, "/") || !strings.Contains(path, "{") {
		return nil
	}

	parts := make([]serveMuxPatternPart, 0, strings.Count(path, "/")+1)
	hasWildcard := false
	for len(path) > 0 {
		path = path[1:]
		if path == "" {
			parts = append(parts, serveMuxPatternPart{kind: serveMuxTrailingSlash})
			break
		}

		separator := strings.IndexByte(path, '/')
		if separator < 0 {
			separator = len(path)
		}
		segment := path[:separator]
		path = path[separator:]

		if !strings.HasPrefix(segment, "{") || len(segment) < 2 || segment[len(segment)-1] != '}' {
			parts = append(parts, serveMuxPatternPart{kind: serveMuxLiteral})
			continue
		}

		name := segment[1 : len(segment)-1]
		if name == "$" {
			parts = append(parts, serveMuxPatternPart{kind: serveMuxEnd})
			break
		}

		name, multi := strings.CutSuffix(name, "...")
		kind := serveMuxWildcard
		if multi {
			kind = serveMuxMultiWildcard
		}
		parts = append(parts, serveMuxPatternPart{kind: kind, name: name})
		hasWildcard = true
	}
	if !hasWildcard {
		return nil
	}
	return parts
}

func setServeMuxPathValues(r *http.Request, parts []serveMuxPatternPart) {
	if len(parts) == 0 {
		return
	}

	path := r.URL.EscapedPath()
	for _, part := range parts {
		switch part.kind {
		case serveMuxLiteral:
			_, path = nextServeMuxSegment(path)
		case serveMuxWildcard:
			value, rest := nextServeMuxSegment(path)
			r.SetPathValue(part.name, serveMuxPathUnescape(value))
			path = rest
		case serveMuxMultiWildcard:
			path = strings.TrimPrefix(path, "/")
			r.SetPathValue(part.name, serveMuxPathUnescape(path))
			return
		default:
			return
		}
	}
}

func nextServeMuxSegment(path string) (string, string) {
	if path == "/" {
		return "/", ""
	}
	path = path[1:]
	separator := strings.IndexByte(path, '/')
	if separator < 0 {
		separator = len(path)
	}
	return path[:separator], path[separator:]
}

func serveMuxPathUnescape(path string) string {
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return path
	}
	return unescaped
}
