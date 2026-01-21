package ngamux

import (
	"encoding/json"
	"log/slog"
)

// Config define ngamux global configuration
type Config struct {
	RemoveTrailingSlash bool
	LogLevel            slog.Level
	JSONMarshal         func(any) ([]byte, error)
	JSONUnmarshal       func([]byte, any) error
}

// NewConfig returns Config with some default values
func NewConfig() Config {
	config := Config{
		RemoveTrailingSlash: true,
		LogLevel:            slog.LevelError,
		JSONMarshal:         json.Marshal,
		JSONUnmarshal:       json.Unmarshal,
	}

	return config
}
