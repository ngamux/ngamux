package ngamux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkNgamuxHotPath benchmarks the ngamux router hot path with various route types.
// This measures the allocation-free routing after the performance fixes.
func BenchmarkNgamuxHotPath(b *testing.B) {
	mux := New(WithLogLevel(LogLevelQuiet))
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
		tc := tc // capture for closure
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

// BenchmarkNgamuxDeepRoutes benchmarks routes with many path segments.
func BenchmarkNgamuxDeepRoutes(b *testing.B) {
	mux := New(WithLogLevel(LogLevelQuiet))
	for i := 0; i < 20; i++ {
		mux.Get("/a/b/c/d/e/f/g/h/i/j", func(w http.ResponseWriter, r *http.Request) {})
	}

	req := httptest.NewRequest(http.MethodGet, "/a/b/c/d/e/f/g/h/i/j", nil)
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkNgamuxAllFallback benchmarks the ALL method fallback behavior.
func BenchmarkNgamuxAllFallback(b *testing.B) {
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.All("/", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkNgamuxExactMethod benchmarks exact method matching.
func BenchmarkNgamuxExactMethod(b *testing.B) {
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/only-get", func(w http.ResponseWriter, r *http.Request) {})

	cases := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, method := range cases {
		b.Run(method, func(b *testing.B) {
			req := httptest.NewRequest(method, "/only-get", nil)
			w := httptest.NewRecorder()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				mux.ServeHTTP(w, req)
			}
		})
	}
}

// BenchmarkNgamux405 benchmarks 405 Method Not Allowed handling.
func BenchmarkNgamux405(b *testing.B) {
	mux := New(WithLogLevel(LogLevelQuiet))
	mux.Get("/get-only", func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodPost, "/get-only", nil)
	w := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}