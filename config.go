package ngamux

// Config define ngamux global configuration
type Config struct {
	RemoveTrailingSlash bool
	LogLevel            LogLevel
}

// NewConfig returns Config with some default values
func NewConfig() Config {
	config := Config{
		RemoveTrailingSlash: true,
		LogLevel:            LogLevelError,
	}

	return config
}
