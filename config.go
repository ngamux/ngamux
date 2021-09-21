package ngamux

type Config struct {
	RemoveTrailingSlash bool
	GlobalErrorHandler  Handler
}

func buildConfig(configs ...Config) Config {
	config := Config{
		RemoveTrailingSlash: true,
	}
	if len(configs) > 0 {
		config = configs[0]
	}
	if config.GlobalErrorHandler == nil {
		config.GlobalErrorHandler = globalErrorHandler
	}

	return config
}
