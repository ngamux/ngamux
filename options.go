package ngamux

func WithTrailingSlash() func(*Config) {
	return func(c *Config) {
		c.RemoveTrailingSlash = false
	}
}

func WithErrorHandler(globalErrorHandler Handler) func(*Config) {
	return func(c *Config) {
		c.GlobalErrorHandler = globalErrorHandler
	}
}
