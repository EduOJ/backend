package logging

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type _logger interface {
	addWriter(writer _writer)
	removeWriter(t reflect.Type)
	Debug(items ...interface{})
	Info(items ...interface{})
	Warning(items ...interface{})
	Error(items ...interface{})
	Fatal(items ...interface{})
	DebugF(format string, items ...interface{})
	InfoF(format string, items ...interface{})
	WarningF(format string, items ...interface{})
	ErrorF(format string, items ...interface{})
	FatalF(format string, items ...interface{})
}

type logger struct {
	writers []_writer //Writers.
}

// Add a writer to the logger.
func (l *logger) addWriter(writer _writer) {
	l.writers = append(l.writers, writer)
}

// Remove all writers of specific type.
// Should not be used. All calls of this function
// must provide a reason in it's comment.
func (l *logger) removeWriter(t reflect.Type) {
	oldWriters := l.writers
	l.writers = make([]_writer, 0)
	for _, w := range oldWriters {
		if reflect.TypeOf(w) != t {
			l.writers = append(l.writers, w)
		}
	}
}

func (l *logger) logWithLevelString(level Level, text string) {
	// Get caller
	var caller string
	pc, _, line, ok := runtime.Caller(2)
	if ok {
		f := runtime.FuncForPC(pc)
		caller = fmt.Sprint(f.Name(), ":", line)
	} else {
		caller = "unknown"
	}
	log := Log{
		Level:  level,
		Time:   time.Now(),
		Text:   text,
		Caller: caller,
	}
	for _, w := range l.writers {
		w.log(log)
	}
}

func (l *logger) logWithLevel(level Level, items ...interface{}) {
	l.logWithLevelString(level, fmt.Sprint(items...))
}

func (l *logger) logWithLevelF(level Level, format string, items ...interface{}) {
	l.logWithLevelString(level, fmt.Sprintf(format, items...))
}

func (l *logger) Debug(items ...interface{}) {
	l.logWithLevel(DEBUG, items...)
}

func (l *logger) Info(items ...interface{}) {
	l.logWithLevel(INFO, items...)
}

func (l *logger) Warning(items ...interface{}) {
	l.logWithLevel(WARNING, items...)
}

func (l *logger) Error(items ...interface{}) {
	l.logWithLevel(ERROR, items...)
}

func (l *logger) Fatal(items ...interface{}) {
	l.logWithLevel(FATAL, items...)
}

func (l *logger) DebugF(format string, items ...interface{}) {
	l.logWithLevelF(DEBUG, format, items...)
}

func (l *logger) InfoF(format string, items ...interface{}) {
	l.logWithLevelF(INFO, format, items...)
}

func (l *logger) WarningF(format string, items ...interface{}) {
	l.logWithLevelF(WARNING, format, items...)
}

func (l *logger) ErrorF(format string, items ...interface{}) {
	l.logWithLevelF(ERROR, format, items...)
}

func (l *logger) FatalF(format string, items ...interface{}) {
	l.logWithLevelF(FATAL, format, items...)
}
