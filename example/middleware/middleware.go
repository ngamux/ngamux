package main

import (
	"fmt"
	"net/http"

	"github.com/ngamux/ngamux"
	"github.com/ngamux/ngamux/middleware"
)

func main() {
	mux := ngamux.NewNgamux()
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "GET /")
	})

	users := mux.Group("/users")
	users.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "GET /users")
	})

	users.Get("/:id", func(rw http.ResponseWriter, r *http.Request) {
		id := ngamux.GetParam(r, "id")
		if id == "123" {
			panic("something error")
		}
		fmt.Fprintf(rw, "GET /users/:id with id = %s", ngamux.GetParam(r, "id"))
	})

	users.Get("/:id/:slug", func(rw http.ResponseWriter, r *http.Request) {
		ngamux.JSON(rw, map[string]string{
			"id":   ngamux.GetParam(r, "id"),
			"slug": ngamux.GetParam(r, "slug"),
		})
	})

	mux.Use(
		middleware.Recover(),
		middleware.EnableCors(middleware.CorsOptions{
			AllowMethods:      []string{http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowHeaders:      []string{"*"},
			AllowOrigins:      []string{"*"},
			AllowExposeHeader: []string{"*"},
		}),
	)

	http.ListenAndServe(":8080", mux)
}
