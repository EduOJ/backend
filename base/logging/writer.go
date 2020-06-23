package logging

import "github.com/leoleoasd/EduOJBackend/base/event"

// Writers should only be used in this package.
// Other packages should use listeners to receive logs.
type _writer interface {
	log(log Log)
	init() error
}

// Writes to the console.
type consoleWriter struct {
	Format string // The format of logger.
}

// Writes to the database for reading from web.
type databaseWriter struct{}

// Calling log listeners.
type eventWriter struct{}

func (w *consoleWriter) log(l Log) {
	// TODO
}

func (w *consoleWriter) init() {}

func (w *databaseWriter) log(l Log) {
	// TODO
}

func (w *databaseWriter) init() {
	// TODO
}

func (w *eventWriter) log(l Log) {
	_, _ = event.FireEvent("on_log", l)
}

func (w *eventWriter) init() {}
