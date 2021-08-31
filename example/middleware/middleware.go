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

	mux.Get("/",
		MiddlewareHello(
			MiddlewareHello(
				MiddlewareHello(
					MiddlewareHello(
						func(rw http.ResponseWriter, r *http.Request) {
							fmt.Fprintln(rw, "hello from handler")
						},
					),
				),
			),
		),
	)

	users := mux.Group("/users")
	users.Get("/",
		MiddlewareHello(
			MiddlewareHello(
				MiddlewareHello(
					MiddlewareHello(
						func(rw http.ResponseWriter, r *http.Request) {
							fmt.Fprintln(rw, "hello from users handler")
						},
					),
				),
			),
		),
	)

	http.ListenAndServe(":8080", mux)
}
