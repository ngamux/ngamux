package ngamux

import "log/slog"

// Config define ngamux global configuration
type Config struct {
	RemoveTrailingSlash bool
	LogLevel            slog.Level
}

// NewConfig returns Config with some default values
func NewConfig() Config {
	config := Config{
		RemoveTrailingSlash: true,
		LogLevel:            slog.LevelError,
	}

	return config
}
