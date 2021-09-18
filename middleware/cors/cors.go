package cors

import (
	"net/http"
	"strings"

	"github.com/ngamux/ngamux"
)

var configDefault = Config{
	AllowOrigins: "*",
	AllowMethods: strings.Join([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}, ","),
	AllowHeaders: "",
}

func New(config ...Config) func(next ngamux.HandlerFunc) ngamux.HandlerFunc {
	cfg := configDefault

	if len(config) > 0 {
		cfg = config[0]
		if cfg.AllowMethods == "" {
			cfg.AllowMethods = configDefault.AllowMethods
		}
		if cfg.AllowOrigins == "" {
			cfg.AllowOrigins = configDefault.AllowOrigins
		}

		cfg.AllowOrigins = strings.ReplaceAll(cfg.AllowOrigins, " ", "")
		cfg.AllowMethods = strings.ReplaceAll(cfg.AllowMethods, " ", "")
		cfg.AllowHeaders = strings.ReplaceAll(cfg.AllowHeaders, " ", "")
	}

	allowOrigins := strings.Split(strings.ReplaceAll(cfg.AllowOrigins, " ", ""), ",")

	return func(next ngamux.HandlerFunc) ngamux.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) error {
			// Every Request
			origin := r.Referer()
			allowOrigin := ""
			if len(origin) > 0 {
				if origin[len(origin)-1] == byte('/') {
					origin = origin[:len(origin)-1]
				}

				for _, ao := range allowOrigins {
					if ao == "*" {
						allowOrigin = ao
						break
					}
					if matchSubdomain(origin, ao) {
						allowOrigin = origin
						break
					}
				}
				rw.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			}

			// Normal Request
			if r.Method != http.MethodOptions {
				return next(rw, r)
			}

			// Preflight Request
			rw.Header().Set("Access-Control-Allow-Methods", cfg.AllowMethods)
			rw.Header().Set("Access-Control-Allow-Headers", cfg.AllowHeaders)
			rw.WriteHeader(http.StatusNoContent)
			return nil
		}
	}
}
