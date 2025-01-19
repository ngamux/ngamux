package ngamux

// WithTrailingSlash returns function that adds RemoveTrailingSlash into config
func WithTrailingSlash() func(*Config) {
	return func(c *Config) {
		c.RemoveTrailingSlash = false
	}
}

// WithLogLevel returns function that adds GlobalErrorHandler into config
func WithLogLevel(level LogLevel) func(*Config) {
	return func(c *Config) {
		c.LogLevel = level
	}
}
