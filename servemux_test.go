package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestServeMuxRequestState(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Get("/users/{id}/posts/{slug}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(r.Pattern + "|" + r.PathValue("id") + "|" + r.PathValue("slug"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/42/posts/hello", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("GET /users/{id}/posts/{slug}|42|hello", rec.Body.String())
}

func TestServeMuxAllWildcardState(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.All("/all/{id}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(r.Pattern + "|" + r.PathValue("id"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/all/42", nil)
	mux.ServeHTTP(rec, req)
	must.Equal(" /all/{id}|42", rec.Body.String())
}

func TestServeMuxHeadFallbackState(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Get("/items/{id}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(r.Pattern + "|" + r.PathValue("id"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, "/items/7", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("GET /items/{id}|7", rec.Body.String())
}

func TestServeMuxMultiWildcardState(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Get("/files/{rest...}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(r.PathValue("rest"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/files/a%2Fb/c", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("a/b/c", rec.Body.String())
}

func TestServeMuxMiddlewareSeesState(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text(r.Pattern + "|" + r.PathValue("id"))
		}
	})
	mux.Get("/middleware/{id}", func(rw http.ResponseWriter, r *http.Request) {})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/middleware/42", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("GET /middleware/{id}|42", rec.Body.String())
}

func TestServeMuxGroupRequestState(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	group := mux.Group("/api")
	group.Get("/users/{id}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text(r.Pattern + "|" + r.PathValue("id"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/users/42", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("GET /api/users/{id}|42", rec.Body.String())
}

func TestServeMuxRedirects(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Get("/dir/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("dir")
	})
	mux.Get("/clean/{id}", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("clean")
	})

	cases := []struct {
		path     string
		location string
	}{
		{"/dir", "/dir/"},
		{"/clean//42", "/clean/42"},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			mux.ServeHTTP(rec, req)
			must.Equal(http.StatusTemporaryRedirect, rec.Code)
			must.Equal(tc.location, rec.Header().Get("Location"))
		})
	}
}

func TestServeMuxNotFoundMiddleware(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	calls := 0
	mux.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			calls++
			if r.URL.Path == "/missing" {
				rw.Header().Set("X-NotFound-Middleware", "1")
			}
			next(rw, r)
		}
	})
	mux.Get("/known", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("known")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/known", nil)
	mux.ServeHTTP(rec, req)
	must.Equal("", rec.Header().Get("X-NotFound-Middleware"))
	must.Equal(1, calls)

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/missing", nil)
	mux.ServeHTTP(rec, req)
	must.Equal(http.StatusNotFound, rec.Code)
	must.Equal("1", rec.Header().Get("X-NotFound-Middleware"))
	must.Equal(2, calls)
}

func TestServeMuxUse(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	middleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
		}
	}
	mux.Use(middleware)
	mux.Use(middleware)
	mux.Use(middleware)

	result := len(mux.middlewares)
	expected := 3

	must.Equal(expected, result)
}

func TestServeMuxGet(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		must := must.New(t)
		mux := NewHttpServeMux()
		mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		must := must.New(t)
		mux := NewHttpServeMux()
		mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/cats", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "404 page not found"
		must.Equal(expected, result)
	})
}

// func TestHead(t *testing.T) {
// 	must := must.New(t)
// 	mux := New(
// 		WithLogLevel(LogLevelQuiet),
// 	)
// 	mux.Head("/", func(rw http.ResponseWriter, r *http.Request) {
// 		Res(rw).Text("ok")
// 	})
//
// 	rec := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodHead, "/", nil)
// 	mux.ServeHTTP(rec, req)
//
// 	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
// 	expected := ""
// 	must.Equal(expected, result)
// }

func TestServeMuxPost(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxPut(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Put("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxPatch(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Patch("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxDelete(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.Delete("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := "ok"
	must.Equal(expected, result)
}

func TestServeMuxAll(t *testing.T) {
	must := must.New(t)
	mux := NewHttpServeMux()
	mux.All("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}
}

func BenchmarkHttpServeMuxHotPath(b *testing.B) {
	noop := func(http.ResponseWriter, *http.Request) {}
	mux := NewHttpServeMux()
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
}

// BenchmarkHttpServeMuxHotPathBenchmark benchmarks the HttpServeMux hot path.
func BenchmarkHttpServeMuxHotPathBenchmark(b *testing.B) {
	mux := NewHttpServeMux()
	mux.Get("/static", func(w http.ResponseWriter, r *http.Request) {})
	mux.Get("/params/{id}/items/{item}", func(w http.ResponseWriter, r *http.Request) {})
	mux.Get("/deep/a/b/c/d/e", func(w http.ResponseWriter, r *http.Request) {})
	mux.Get("/only-get", func(w http.ResponseWriter, r *http.Request) {})

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
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, req)
			}
		})
	}
}

// BenchmarkHttpServeMuxVsStdlibBenchmark compares HttpServeMux against stdlib ServeMux.
func BenchmarkHttpServeMuxVsStdlibBenchmark(b *testing.B) {
	hMux := NewHttpServeMux()
	hMux.Get("/static", func(w http.ResponseWriter, r *http.Request) {})
	hMux.Get("/params/{id}/items/{item}", func(w http.ResponseWriter, r *http.Request) {})
	hMux.Get("/deep/a/b/c/d/e", func(w http.ResponseWriter, r *http.Request) {})
	hMux.Get("/only-get", func(w http.ResponseWriter, r *http.Request) {})

	stdMux := http.NewServeMux()
	stdMux.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {})
	stdMux.HandleFunc("/params/{id}/items/{item}", func(w http.ResponseWriter, r *http.Request) {})
	stdMux.HandleFunc("/deep/a/b/c/d/e", func(w http.ResponseWriter, r *http.Request) {})
	stdMux.HandleFunc("/only-get", func(w http.ResponseWriter, r *http.Request) {})

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
		tc := tc
		b.Run(tc.name+"-http", func(b *testing.B) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				hMux.ServeHTTP(w, req)
			}
		})

		b.Run(tc.name+"-stdlib", func(b *testing.B) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stdMux.ServeHTTP(w, req)
			}
		})
	}
}

// BenchmarkHttpServeMuxDeepRoutesBenchmark benchmarks with many segments.
func BenchmarkHttpServeMuxDeepRoutesBenchmark(b *testing.B) {
	mux := NewHttpServeMux()
	for i := 0; i < 1; i++ {
		mux.Get("/deep/route", func(w http.ResponseWriter, r *http.Request) {})
	}

	req := httptest.NewRequest(http.MethodGet, "/a/b/c/d/e/f/g/h/i/j", nil)
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkHttpServeMuxAllFallbackBenchmark benchmarks ALL method.
func BenchmarkHttpServeMuxAllFallbackBenchmark(b *testing.B) {
	mux := NewHttpServeMux()
	mux.All("/", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

