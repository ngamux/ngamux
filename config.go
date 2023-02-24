package ngamux

// Config define ngamux global configuration
type Config struct {
	RemoveTrailingSlash bool
	GlobalErrorHandler  Handler
}

// NewConfig returns Config with some default values
func NewConfig() Config {
	config := Config{
		RemoveTrailingSlash: true,
		GlobalErrorHandler:  globalErrorHandler,
	}

	return config
}
