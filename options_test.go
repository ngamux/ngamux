package ngamux

import (
	"testing"

	"github.com/golang-must/must"
)

func TestOptions(t *testing.T) {
	t.Run("set RemoveTrailingSlash", func(t *testing.T) {
		must := must.New(t)

		mux := New(WithTrailingSlash())
		must.False(mux.config.RemoveTrailingSlash)
	})

	t.Run("set LogLevel", func(t *testing.T) {
		must := must.New(t)

		mux := New(WithLogLevel(LogLevelQuiet))
		must.Equal(mux.config.LogLevel, LogLevelQuiet)
	})

}
