package main

import (
	"fmt"
	"net/http"

	"github.com/ngamux/ngamux"
)

func MiddlewareHello(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "hello from middleware")
		handler(rw, r)
	}
}

func main() {
	mux := ngamux.NewNgamux(ngamux.Config{
		RemoveTrailingSlash: true,
	})

	mux.Use(func(next ngamux.HandlerFunc) ngamux.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) error {
			fmt.Println("hello from middleware 1")
			return next(rw, r)
		}
	})

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
		fmt.Fprintln(rw, "hello from users handler")
		fmt.Println("hello from handler")
		return nil
	})

	http.ListenAndServe(":8080", mux)
}
