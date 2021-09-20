package ngamux

type Config struct {
	RemoveTrailingSlash bool
	NotFoundHandler     HandlerFunc
}

func buildConfig(configs ...Config) Config {
	config := Config{
		RemoveTrailingSlash: true,
	}
	if len(configs) > 0 {
		config = configs[0]
	}
	if config.NotFoundHandler == nil {
		config.NotFoundHandler = handlerNotFound
	}

	return config
}
