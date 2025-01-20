package ngamux

import (
	"testing"

	"github.com/golang-must/must"
)

func TestLogConstants(t *testing.T) {
	must.Equal(t, LogLevelQuiet.String(), "QUIET")
	must.Equal(t, LogLevelInfo.String(), "INFO")
	must.Equal(t, LogLevelWarn.String(), "WARN")
	must.Equal(t, LogLevelError.String(), "ERRO")
	must.Equal(t, LogLevel(-1).String(), "")
}

func TestIsLogCanShow(t *testing.T) {
	t.Run("quiet", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(LogLevelQuiet))

		must.False(m.isLogCanShow(LogLevelInfo))
		must.False(m.isLogCanShow(LogLevelWarn))
		must.False(m.isLogCanShow(LogLevelError))
	})

	t.Run("info", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(LogLevelInfo))

		must.True(m.isLogCanShow(LogLevelInfo))
		must.False(m.isLogCanShow(LogLevelWarn))
		must.False(m.isLogCanShow(LogLevelError))
	})

	t.Run("warn", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(LogLevelWarn))

		must.True(m.isLogCanShow(LogLevelInfo))
		must.True(m.isLogCanShow(LogLevelWarn))
		must.False(m.isLogCanShow(LogLevelError))
	})

	t.Run("error", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(LogLevelError))

		must.True(m.isLogCanShow(LogLevelInfo))
		must.True(m.isLogCanShow(LogLevelWarn))
		must.True(m.isLogCanShow(LogLevelError))
	})

}
