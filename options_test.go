package ngamux

import (
	"testing"

	"github.com/golang-must/must"
)

func TestWithTrailingSlash(t *testing.T) {
	t.Run("can set false to config.RemoveTrailingSlash", func(t *testing.T) {
		must := must.New(t)

		mux := New(WithTrailingSlash())
		must.False(mux.config.RemoveTrailingSlash)
	})

	t.Run("can set globalErrorHandler to config.GlobalErrorHandler", func(t *testing.T) {
		must := must.New(t)

		mux := New(WithErrorHandler(globalErrorHandler))
		must.NotNil(mux.config.GlobalErrorHandler)
	})
}
