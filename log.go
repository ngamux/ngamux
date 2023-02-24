package ngamux

import (
	"fmt"
)

type LogLevel int

const (
	LogLevelQuiet LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (l LogLevel) String() string {
	switch l {
	case 0:
		return "QUIET"
	case 1:
		return "INFO"
	case 2:
		return "WARN"
	case 3:
		return "ERRO"
	}

	return ""
}

func (m Ngamux) isLogCanShow(level LogLevel) bool {
	if m.config.LogLevel == LogLevelQuiet {
		return false
	}

	if m.config.LogLevel == LogLevelInfo && level == LogLevelInfo {
		return true
	}

	if (m.config.LogLevel == LogLevelWarn && level == LogLevelInfo) || (m.config.LogLevel == LogLevelWarn && level == LogLevelWarn) {
		return true
	}

	if (m.config.LogLevel == LogLevelError && level == LogLevelInfo) || (m.config.LogLevel == LogLevelError && level == LogLevelWarn) || (m.config.LogLevel == LogLevelError && level == LogLevelError) {
		return true
	}

	return false
}

func (m Ngamux) Log(level LogLevel, message string, data ...any) {
	if !m.isLogCanShow(level) {
		return
	}

	fmt.Printf(fmt.Sprintf("[%s] %s\n", level, message), data...)
}
