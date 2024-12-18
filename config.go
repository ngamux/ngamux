package ngamux

import "net/http"

// Config define ngamux global configuration
type Config struct {
	RemoveTrailingSlash bool
	GlobalErrorHandler  http.HandlerFunc
	LogLevel            LogLevel
}

// NewConfig returns Config with some default values
func NewConfig() Config {
	config := Config{
		RemoveTrailingSlash: true,
		GlobalErrorHandler:  globalErrorHandler,
		LogLevel:            LogLevelError,
	}

	return config
}
