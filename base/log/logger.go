package log

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type _logger interface {
	addWriter(writer _writer)
	removeWriter(t reflect.Type)
	isReady() bool
	setReady()
	Debug(items ...interface{})
	Info(items ...interface{})
	Warning(items ...interface{})
	Error(items ...interface{})
	Fatal(items ...interface{})
	Debugf(format string, items ...interface{})
	Infof(format string, items ...interface{})
	Warningf(format string, items ...interface{})
	Errorf(format string, items ...interface{})
	Fatalf(format string, items ...interface{})
}

type logger struct {
	writers []_writer //Writers.
	ready   bool
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

func (l *logger) isReady() bool {
	return l.ready
}

func (l *logger) setReady() {
	l.ready = true
}

func (l *logger) logWithLevelString(level Level, message string) {
	caller := "unknown"
	{
		// Find caller out of the log package.
		pc := make([]uintptr, 5)
		runtime.Callers(1, pc)
		frames := runtime.CallersFrames(pc)
		more := true
		for more {
			var frame runtime.Frame
			frame, more = frames.Next()
			if !strings.Contains(frame.File, "github.com/leoleoasd/EduOJBackend/base/log") {
				caller = fmt.Sprint(frame.Func.Name(), ":", frame.Line)
			}
		}
	}
	log := Log{
		Level:   level,
		Time:    time.Now(),
		Message: message,
		Caller:  caller,
	}
	if l.ready == false {
		// Logger don't been initialized yet.
		// So we should just output to stdout.
		(&consoleWriter{}).log(log)
		return
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

func (l *logger) Debugf(format string, items ...interface{}) {
	l.logWithLevelF(DEBUG, format, items...)
}

func (l *logger) Infof(format string, items ...interface{}) {
	l.logWithLevelF(INFO, format, items...)
}

func (l *logger) Warningf(format string, items ...interface{}) {
	l.logWithLevelF(WARNING, format, items...)
}

func (l *logger) Errorf(format string, items ...interface{}) {
	l.logWithLevelF(ERROR, format, items...)
}

func (l *logger) Fatalf(format string, items ...interface{}) {
	l.logWithLevelF(FATAL, format, items...)
}
