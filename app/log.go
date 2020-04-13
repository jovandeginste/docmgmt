package app

import "log"

const (
	LOG_UNDEFINED = iota
	LOG_CRITICAL
	LOG_ERROR
	LOG_WARNING
	LOG_INFO
	LOG_DEBUG
)

func (a *App) Log(level int, message ...string) {
	if level <= a.Configuration.LogLevel {
		log.Println(message)
	}
}

func (a *App) Logf(level int, message string, args ...interface{}) {
	if level <= a.Configuration.LogLevel {
		log.Printf(message, args...)
	}
}
