package ngamux

// Config define ngamux global configuration
type Config struct {
	RemoveTrailingSlash bool
	GlobalErrorHandler  Handler
}

func NewConfig() Config {
	config := Config{
		RemoveTrailingSlash: true,
		GlobalErrorHandler:  globalErrorHandler,
	}

	return config
}
