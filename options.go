package ngamux

import "net/http"

// WithTrailingSlash returns function that adds RemoveTrailingSlash into config
func WithTrailingSlash() func(*Config) {
	return func(c *Config) {
		c.RemoveTrailingSlash = false
	}
}

// WithErrorHandler returns function that adds GlobalErrorHandler into config
func WithErrorHandler(globalErrorHandler http.HandlerFunc) func(*Config) {
	return func(c *Config) {
		c.GlobalErrorHandler = globalErrorHandler
	}
}

// WithLogLevel returns function that adds GlobalErrorHandler into config
func WithLogLevel(level LogLevel) func(*Config) {
	return func(c *Config) {
		c.LogLevel = level
	}
}
