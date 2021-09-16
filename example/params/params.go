package main

import (
	"fmt"
	"net/http"

	"github.com/ngamux/ngamux"
)

func main() {
	mux := ngamux.NewNgamux()
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
		fmt.Fprintln(rw, "GET /")
		return nil
	})

	users := mux.Group("/users")
	users.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
		fmt.Fprintln(rw, "GET /users")
		return nil
	})

	users.Get("/:id", func(rw http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(rw, "GET /users/:id with id = %s", ngamux.GetParam(r, "id"))
		return nil
	})

	users.Get("/:id/:slug", func(rw http.ResponseWriter, r *http.Request) error {
		return ngamux.JSON(rw, map[string]string{
			"id":   ngamux.GetParam(r, "id"),
			"slug": ngamux.GetParam(r, "slug"),
		})
	})

	http.ListenAndServe(":8080", mux)
}
