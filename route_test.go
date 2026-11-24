package ngamux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestRouteMethodFallback(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/resource", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("get")
	})
	mux.Post("/resource", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("post")
	})
	mux.All("/resource", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("all")
	})
	mux.All("/fallback/{id}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(r.Method + "|" + Req(r).PathValue("id"))
	})

	cases := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/resource", "get"},
		{http.MethodPost, "/resource", "post"},
		{http.MethodDelete, "/resource", "all"},
		{http.MethodPatch, "/fallback/42", "PATCH|42"},
		{http.MethodGet, "/missing", "404 page not found"},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			mux.ServeHTTP(rec, req)
			must.Equal(tc.body, strings.ReplaceAll(rec.Body.String(), "\n", ""))
			must.Equal(tc.method, req.Method)
		})
	}
}

func TestRouteWrongMethod(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/only-get", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("get")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/only-get", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("405 method not allowed", strings.ReplaceAll(rec.Body.String(), "\n", ""))
	must.Equal(http.MethodPost, req.Method)
}

func TestRouteExactChildPrecedesWildcard(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/items/{id}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("wildcard|" + Req(r).PathValue("id"))
	})
	mux.Get("/items/special", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("static")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items/special", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("static", rec.Body.String())

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/items/other", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("wildcard|other", rec.Body.String())
}

func TestRouteFailedMatchDoesNotSetParams(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/users/{id}/posts/{slug}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("matched")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42/missing/hello", nil)
	handler, pattern := mux.Handler(req)
	must.Nil(handler)
	must.Equal("", pattern)
	must.Equal("", req.PathValue("id"))
}

func TestRouteParamOverflow(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))

	const count = maxStackRouteParams + 2
	var route strings.Builder
	var path strings.Builder
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("p%d", i)
		value := fmt.Sprintf("v%d", i)
		route.WriteString("/{")
		route.WriteString(name)
		route.WriteString("}")
		path.WriteString("/")
		path.WriteString(value)
	}
	mux.Get(route.String(), func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(Req(r).PathValue("p0") + "|" + Req(r).PathValue(fmt.Sprintf("p%d", count-1)))
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path.String(), nil)
	mux.ServeHTTP(rec, req)
	must.Equal("v0|v9", rec.Body.String())
}

func TestRoutePathWithSpaces(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/files/{name}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(Req(r).PathValue("name"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL.Path = "/files/a b"
	mux.ServeHTTP(rec, req)
	must.Equal("a b", rec.Body.String())
}

func TestGetRoute(t *testing.T) {
	must := must.New(t)
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler, pattern := mux.Handler(req)
	must.NotNil(handler)
	must.Equal("/", pattern)
	handler.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)

	mux1 := New(WithLogLevel(LogLevelQuiet))
	mux1.Get("/{a}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	req1 := httptest.NewRequest(http.MethodGet, "/123", nil)
	rec1 := httptest.NewRecorder()
	handler1, pattern1 := mux1.Handler(req1)
	must.NotNil(handler1)
	must.Equal("/{a}", pattern1)
	handler1.ServeHTTP(rec1, req1)

	result = strings.ReplaceAll(rec1.Body.String(), "\n", "")
	expected = "ok"
	must.Equal(expected, result)

	req2 := httptest.NewRequest(http.MethodGet, "/123", nil)
	handler2, pattern2 := mux.Handler(req2)
	must.Nil(handler2)
	must.Equal("", pattern2)

	req3 := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec3 := httptest.NewRecorder()
	handler3, pattern3 := mux.Handler(req3)
	must.NotNil(handler3)
	must.Equal("/", pattern3)
	handler3.ServeHTTP(rec3, req3)
	result3 := strings.ReplaceAll(rec3.Body.String(), "\n", "")
	expected3 := "405 method not allowed"
	must.Equal(expected3, result3)
}

type discardResponseWriter struct {
	header http.Header
}

func newDiscardResponseWriter() *discardResponseWriter {
	return &discardResponseWriter{header: make(http.Header)}
}

func (w *discardResponseWriter) Header() http.Header               { return w.header }
func (w *discardResponseWriter) Write(p []byte) (int, error)       { return len(p), nil }
func (w *discardResponseWriter) WriteString(s string) (int, error) { return len(s), nil }
func (w *discardResponseWriter) WriteHeader(int)                   {}

func TestRouterHotPathAllocations(t *testing.T) {
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/static", func(http.ResponseWriter, *http.Request) {})
	mux.Get("/params/{id}/items/{item}", func(http.ResponseWriter, *http.Request) {})

	cases := []struct {
		name string
		path string
	}{
		{"static", "/static"},
		{"params", "/params/123/items/abc"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			w := newDiscardResponseWriter()
			allocs := testing.AllocsPerRun(100, func() {
				mux.ServeHTTP(w, req)
			})
			if allocs != 0 {
				t.Fatalf("router allocated %v objects per request, want 0", allocs)
			}
		})
	}
}

func BenchmarkRouterHotPath(b *testing.B) {
	noop := func(http.ResponseWriter, *http.Request) {}
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/static", noop)
	mux.Get("/params/{id}/items/{item}", noop)
	mux.Get("/deep/a/b/c/d/e", noop)
	mux.Get("/only-get", noop)

	cases := []struct {
		name   string
		method string
		path   string
	}{
		{"static", http.MethodGet, "/static"},
		{"params", http.MethodGet, "/params/123/items/abc"},
		{"deep", http.MethodGet, "/deep/a/b/c/d/e"},
		{"miss", http.MethodGet, "/missing"},
		{"wrong_method", http.MethodPost, "/only-get"},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := newDiscardResponseWriter()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				mux.ServeHTTP(w, req)
			}
		})
	}

	b.Run("all_fallback", func(b *testing.B) {
		fallback := New(WithLogLevel(LogLevelQuiet))
		fallback.All("/all/{id}", noop)
		req := httptest.NewRequest(http.MethodPatch, "/all/42", nil)
		w := newDiscardResponseWriter()
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			fallback.ServeHTTP(w, req)
		}
	})
}
