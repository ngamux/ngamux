package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ngamux/ngamux"
)

func customPanicHandler(w http.ResponseWriter, r *http.Request, p interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Oops, my bad")
	log.Printf("panic: %v on: %v", p, r.URL.Path)
}

func main() {
	mux := ngamux.NewNgamux(ngamux.Config{
		//PanicHandler: ngamux.DefaultPanicHandler,
		PanicHandler: customPanicHandler,
	})

	mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "not panicking")
	})

	mux.Get("/panic", func(rw http.ResponseWriter, r *http.Request) {
		panic("panic testing")
	})

	http.ListenAndServe(":8080", mux)
}
