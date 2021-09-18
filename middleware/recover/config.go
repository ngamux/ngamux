package recover

import "net/http"

type Config struct {
	ErrorHandler func(http.ResponseWriter, *http.Request, error)
}
