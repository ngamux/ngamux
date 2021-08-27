package main

import (
	"fmt"
	"net/http"

	"github.com/ngamux/ngamux"
)

func main() {
	mux := ngamux.NewNgamux()
	users := mux.Group("/users")
	users.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "GET /users")
	})

	http.ListenAndServe(":8080", mux)
}
