package main

import (
	"fmt"
	"net/http"

	"github.com/ngamux/ngamux"
)

func main() {
	mux := ngamux.NewNgamux()
	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "GET /")
	})

	mux.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		in := map[string]string{}
		err := ngamux.GetBody(r, &in)
		if err != nil {
			ngamux.JSONWithStatus(rw, http.StatusBadRequest, err.Error())
		}
		ngamux.JSON(rw, in)
	})

	http.ListenAndServe(":8080", mux)
}
