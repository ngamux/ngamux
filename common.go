package ngamux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SETUP
func WithMiddlewares(middleware ...MiddlewareFunc) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		h := next
		for i := len(middleware) - 1; i >= 0; i-- {
			h = middleware[i](h)
		}
		return h
	}
}

// REQUEST
func GetParam(r *http.Request, key string) string {
	params := r.Context().Value(KeyContextParams).([][]string)
	for _, param := range params {
		if param[0] == key {
			return param[1]
		}
	}

	return ""
}

func GetBody(r *http.Request, store interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		return err
	}

	return nil
}

func ToContext(r *http.Request, key interface{}, value interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	r = r.WithContext(ctx)
	return r
}

func FromContext(r *http.Request, key interface{}) interface{} {
	value := r.Context().Value(key)
	return value
}

// RESPONSE
func JSON(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Add("content-type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprint(rw, string(jsonData))
	return nil
}

func JSONWithStatus(rw http.ResponseWriter, status int, data interface{}) error {
	rw.WriteHeader(status)
	err := JSON(rw, data)
	if err != nil {
		return err
	}

	return nil
}
