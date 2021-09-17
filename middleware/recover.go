package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ngamux/ngamux"
)

type recoverHandler func(http.ResponseWriter, *http.Request, error)

func Recover(onerror ...recoverHandler) func(next ngamux.HandlerFunc) ngamux.HandlerFunc {
	if len(onerror) < 1 {
		onerror = append(onerror, func(rw http.ResponseWriter, r *http.Request, e error) {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(rw, e)
			log.Println("error:", e)
		})
	}

	return func(next ngamux.HandlerFunc) ngamux.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) error {
			defer func() {
				if err := recover(); err != nil {
					onerror[0](rw, r, errors.New(err.(string)))
				}
			}()
			return next(rw, r)
		}
	}
}
