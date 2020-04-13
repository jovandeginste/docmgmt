package app

import "log"

const (
	LogUndefined = iota
	LogCritical
	LogError
	LogWarning
	LogInfo
	LogDebug
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
