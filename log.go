package ngamux

import (
	"context"
	"fmt"
	"log/slog"
)

const LogLevelQuiet slog.Level = -8

func (m Ngamux) isLogCanShow(level slog.Level) bool {
	if m.config.LogLevel == LogLevelQuiet {
		return false
	}

	if m.config.LogLevel == slog.LevelInfo && level == slog.LevelInfo {
		return true
	}

	if (m.config.LogLevel == slog.LevelWarn && level == slog.LevelInfo) ||
		(m.config.LogLevel == slog.LevelWarn && level == slog.LevelWarn) {
		return true
	}

	if (m.config.LogLevel == slog.LevelError && level == slog.LevelInfo) ||
		(m.config.LogLevel == slog.LevelError && level == slog.LevelWarn) ||
		(m.config.LogLevel == slog.LevelError && level == slog.LevelError) {
		return true
	}

	return false
}

func (m Ngamux) Log(level slog.Level, message string, data ...any) {
	if !m.isLogCanShow(level) {
		return
	}

	slog.Default().Log(context.Background(), level, fmt.Sprintf("[%s] %s\n", level, message), data...)
}
