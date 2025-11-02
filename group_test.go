package ngamux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-must/must"
)

func TestGroup(t *testing.T) {
	must := must.New(t)
	mux := New(
		WithLogLevel(LogLevelQuiet),
	)
	handler := func(rw http.ResponseWriter, r *http.Request) {
		Res(rw).Text("ok")
	}

	{
		a := mux.Group("/a")
		a.Get("", handler)
		a.Post("", handler)
		a.Put("", handler)
		a.Patch("", handler)
		a.Delete("", handler)
		a.All("", handler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/a", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}

	{
		a := mux.Group("/a")
		b := a.Group("/b")
		b.Get("", handler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/a/b", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	}
}

func TestServeMuxGroup(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		must := must.New(t)
		mux := NewHttpServeMux()
		handler := func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		}
		a := mux.Group("/a")
		a.Get("/", handler)
		a.Post("/", handler)
		a.Put("", handler)
		a.Patch("", handler)
		a.Delete("", handler)
		a.All("", handler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/a", nil)
		mux.ServeHTTP(rec, req)

		result := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expected := "ok"
		must.Equal(expected, result)
	})

	var add = func(result *[]int, n int) MiddlewareFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				*result = append(*result, n)
				next(w, r)
			}
		}
	}

	t.Run("middlewares global", func(t *testing.T) {
		must := must.New(t)
		result := []int{}
		expected := []int{1, 2}
		mux := NewHttpServeMux()

		mux.Use(add(&result, 1), add(&result, 2))

		handler := func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		}
		mux.Get("/", handler)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		mux.ServeHTTP(rec, req)

		resultHttp := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expectedHttp := "ok"
		must.Equal(expectedHttp, resultHttp)
		must.Equal(expected, result)
	})

	t.Run("middlewares global + group", func(t *testing.T) {
		must := must.New(t)
		result := []int{}
		expected := []int{1, 2, 3, 4}
		mux := NewHttpServeMux()

		mux.Use(add(&result, 1), add(&result, 2))
		v2 := mux.Group("/v2")
		v2.Use(add(&result, 3), add(&result, 4))

		handler := func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		}
		v2.Get("", handler)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v2", nil)
		mux.ServeHTTP(rec, req)

		resultHttp := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expectedHttp := "ok"
		must.Equal(expectedHttp, resultHttp)
		must.Equal(expected, result)
	})

	t.Run("middlewares global + group nested", func(t *testing.T) {
		must := must.New(t)
		result := []int{}
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		mux := NewHttpServeMux()

		mux.Use(add(&result, 1), add(&result, 2))
		v2 := mux.Group("/v2")
		v2.Use(add(&result, 3), add(&result, 4))

		users := v2.Group("/users")
		users.Use(add(&result, 5), add(&result, 6))

		handler := func(rw http.ResponseWriter, r *http.Request) {
			Res(rw).Text("ok")
		}
		users.GroupFunc("/{id}", func(mux *HttpServeMux) {
			users.Use(add(&result, 7), add(&result, 8))
			mux.Get("/txs", handler, add(&result, 10), add(&result, 9))
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v2/users/1/txs", nil)
		mux.ServeHTTP(rec, req)

		resultHttp := strings.ReplaceAll(rec.Body.String(), "\n", "")
		expectedHttp := "ok"
		must.Equal(expectedHttp, resultHttp)
		must.Equal(expected, result)
	})
}
