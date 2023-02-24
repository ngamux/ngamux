package ngamux

import (
	"net/http"
	"testing"

	"github.com/golang-must/must"
)

func TestWithMiddlewares(t *testing.T) {
	must := must.New(t)
	result := WithMiddlewares()(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	must.NotNil(result)

	result = WithMiddlewares(nil)(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	must.NotNil(result)

	result = WithMiddlewares(nil)(nil)
	must.Nil(result)
}
