package recover

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ngamux/ngamux"
)

var configDefault = Config{
	ErrorHandler: func(rw http.ResponseWriter, r *http.Request, e error) {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(rw, e)
		log.Println("error:", e)
	},
}

func New(config ...Config) func(next ngamux.HandlerFunc) ngamux.HandlerFunc {
	cfg := configDefault

	if len(config) > 0 {
		cfg = config[0]
	}

	return func(next ngamux.HandlerFunc) ngamux.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) error {
			defer func() {
				if err := recover(); err != nil {
					cfg.ErrorHandler(rw, r, errors.New(err.(string)))
				}
			}()
			return next(rw, r)
		}
	}
}
