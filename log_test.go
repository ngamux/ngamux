package ngamux

import (
	"bytes"
	"log"
	"log/slog"
	"testing"

	"github.com/golang-must/must"
)

func TestLogConstants(t *testing.T) {
	must.Equal(t, LogLevelQuiet.String(), "DEBUG-4")
}

func TestIsLogCanShow(t *testing.T) {
	t.Run("quiet", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(LogLevelQuiet))

		must.False(m.isLogCanShow(slog.LevelInfo))
		must.False(m.isLogCanShow(slog.LevelWarn))
		must.False(m.isLogCanShow(slog.LevelError))
	})

	t.Run("info", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(slog.LevelInfo))

		must.True(m.isLogCanShow(slog.LevelInfo))
		must.False(m.isLogCanShow(slog.LevelWarn))
		must.False(m.isLogCanShow(slog.LevelError))
	})

	t.Run("warn", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(slog.LevelWarn))

		must.True(m.isLogCanShow(slog.LevelInfo))
		must.True(m.isLogCanShow(slog.LevelWarn))
		must.False(m.isLogCanShow(slog.LevelError))
	})

	t.Run("error", func(t *testing.T) {
		must := must.New(t)
		m := New(WithLogLevel(slog.LevelError))

		must.True(m.isLogCanShow(slog.LevelInfo))
		must.True(m.isLogCanShow(slog.LevelWarn))
		must.True(m.isLogCanShow(slog.LevelError))
	})

}

func TestLog(t *testing.T) {
	mux := New()
	b := &bytes.Buffer{}
	log.SetOutput(b)
	mux.Log(LogLevelQuiet, "ok")
	must.Equal(t, b.String(), "")

	mux.Log(slog.LevelInfo, "")
	must.NotEqual(t, b.String(), "")
}
