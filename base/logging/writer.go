package logging

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base/event"
)

// Writers should only be used in this package.
// Other packages should use listeners to receive logs.
type _writer interface {
	log(log Log)
}

// Writes to the console.
type consoleWriter struct {
}

// Writes to the database for reading from web.
type databaseWriter struct{}

// Calling log listeners.
type eventWriter struct{}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

var (
	colors = []string{
		FATAL:   colorSeq(colorMagenta),
		ERROR:   colorSeq(colorRed),
		WARNING: colorSeq(colorYellow),
		INFO:    colorSeq(colorGreen),
		DEBUG:   colorSeq(colorCyan),
	}
)

func colorSeq(color int) string {
	return fmt.Sprintf("\033[%dm", color)
}

func (w *consoleWriter) log(l Log) {
	fmt.Printf("%s[%s][%s] ▶ %s\u001B[0m %s\n",
		colors[l.Level],
		l.Time.Format("15:04:05"),
		l.Caller,
		l.Level.String(),
		l.Message,
	)
}

func (w *databaseWriter) log(l Log) {
	// TODO
}

func (w *databaseWriter) init() (err error) {
	// TODO
	return
}

func (w *eventWriter) log(l Log) {
	event.FireEvent("log", l)
}
