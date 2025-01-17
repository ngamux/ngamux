package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestNewNgamux(t *testing.T) {
	must := must.New(t)
	result := New(
		WithLogLevel(LogLevelQuiet),
	)
	config := NewConfig()
	expected := &Ngamux{
		config: &config,
	}

	must.Equal(expected.config.RemoveTrailingSlash, result.config.RemoveTrailingSlash)
}

func TestUse(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
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

func TestConfig(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)

	result := mux.Config()
	must.Equal(result.RemoveTrailingSlash, true)
	must.Equal(result.LogLevel, LogLevelQuiet)
}

func TestGet(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		must := must.New(t)
		mux := New(
			WithLogLevel(LogLevelQuiet),
		)
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

}

func TestHead(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Head("/", func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, "/", nil)
	mux.ServeHTTP(rec, req)

	result := strings.ReplaceAll(rec.Body.String(), "\n", "")
	expected := ""
	must.Equal(expected, result)
}

func TestPost(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
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

func TestPut(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
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

func TestPatch(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
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

func TestDelete(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
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

func TestAll(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
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

func TestHandle(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	}))

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}

	{
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/cats", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "404 page not found"
		must.Equal(expected, result)
	}
}

func TestWith(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)

	{
		mux1 := mux.With(func(next http.HandlerFunc) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				next(rw, r)
			}
		}, func(next http.HandlerFunc) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				next(rw, r)
			}
		})
		must.NotNil(mux1)
		must.NotNil(mux1.parent)
	}
}

func BenchmarkNgamux(b *testing.B) {
	h1 := func(w http.ResponseWriter, r *http.Request) {}
	h2 := func(w http.ResponseWriter, r *http.Request) {}
	h3 := func(w http.ResponseWriter, r *http.Request) {}
	h4 := func(w http.ResponseWriter, r *http.Request) {}
	h5 := func(w http.ResponseWriter, r *http.Request) {}
	h6 := func(w http.ResponseWriter, r *http.Request) {}

	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	mux.Get("/", h1)
	mux.Get("/general", h2)
	mux.Get("/general/:id/and/:this", h3)

	mux1 := mux.Group("/group/:x/:hash")
	mux1.Get("/", h4)          // subrouter-1
	mux1.Get("/{network}", h5) // subrouter-1
	mux1.Get("/twitter", h5)

	mux2 := mux.Group("/direct")
	mux2.Get("/", h6) // subrouter-2
	mux2.Get("/download", h6)

	routes := []string{
		"/",
		"/general",
		"/general/123/and/this",
		"/general/123/foo/this",
		"/group/z/aBc",                 // subrouter-1
		"/group/z/aBc/twitter",         // subrouter-1
		"/group/z/aBc/direct",          // subrouter-2
		"/group/z/aBc/direct/download", // subrouter-2
	}

	for _, path := range routes {
		b.Run("route/"+path, func(b *testing.B) {
			b.ResetTimer()
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}

func BenchmarkHttpServeMux(b *testing.B) {
	h1 := func(w http.ResponseWriter, r *http.Request) {}
	h2 := func(w http.ResponseWriter, r *http.Request) {}
	h3 := func(w http.ResponseWriter, r *http.Request) {}
	h4 := func(w http.ResponseWriter, r *http.Request) {}
	h5 := func(w http.ResponseWriter, r *http.Request) {}
	h6 := func(w http.ResponseWriter, r *http.Request) {}

	mux := NewHttpServeMux()
	mux.Get("/", h1)
	mux.Get("/general", h2)
	mux.Get("/general/:id/and/:this", h3)

	mux1 := mux.Group("/group/:x/:hash")
	mux1.Get("/", h4)          // subrouter-1
	mux1.Get("/{network}", h5) // subrouter-1
	mux1.Get("/twitter", h5)

	mux2 := mux.Group("/direct")
	mux2.Get("/", h6) // subrouter-2
	mux2.Get("/download", h6)

	routes := []string{
		"/",
		"/general",
		"/general/123/and/this",
		"/general/123/foo/this",
		"/group/z/aBc",                 // subrouter-1
		"/group/z/aBc/twitter",         // subrouter-1
		"/group/z/aBc/direct",          // subrouter-2
		"/group/z/aBc/direct/download", // subrouter-2
	}

	for _, path := range routes {
		b.Run("route/"+path, func(b *testing.B) {
			b.ResetTimer()
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, r)
			}
		})
	}
}
