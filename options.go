package ngamux

// WithTrailingSlash returns function that adds RemoveTrailingSlash into config
func WithTrailingSlash() func(*Config) {
	return func(c *Config) {
		c.RemoveTrailingSlash = false
	}
}

// WithErrorHandler returns function that adds GlobalErrorHandler into config
func WithErrorHandler(globalErrorHandler Handler) func(*Config) {
	return func(c *Config) {
		c.GlobalErrorHandler = globalErrorHandler
	}
}
