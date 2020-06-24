package logging

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

type fakeWriter struct {
	lastLog Log
}

func (f *fakeWriter) init() error {
	return nil
}
func (f *fakeWriter) log(l Log) {
	f.lastLog = l
}

func TestLogger(t *testing.T) {
	var l _logger
	l = &logger{}
	w := &fakeWriter{}
	l.addWriter(w)
	assert.Equal(t, w, l.(*logger).writers[0], "Writer should be fake writer.")

	levels := []struct {
		l            Level
		logFunction  func(items ...interface{})
		logFFunction func(format string, items ...interface{})
	}{
		{
			DEBUG,
			l.Debug,
			l.Debugf,
		}, {
			INFO,
			l.Info,
			l.Infof,
		}, {
			WARNING,
			l.Warning,
			l.Warningf,
		}, {
			ERROR,
			l.Error,
			l.Errorf,
		}, {
			FATAL,
			l.Fatal,
			l.Fatalf,
		},
	}
	for _, level := range levels {
		level.logFunction(123, "321")
		assert.Equal(t, w.lastLog.Level, level.l, "Level should be same.")
		assert.Less(t, time.Since(w.lastLog.Time).Nanoseconds(), time.Millisecond.Nanoseconds(), "Time difference should be less than 1 ms.")
		_, _, line, _ := runtime.Caller(0)
		assert.Equal(t, w.lastLog.Caller, fmt.Sprint("github.com/leoleoasd/EduOJBackend/base/logging.TestLogger:", line-3), "Caller should be this function.")
		assert.Equal(t, w.lastLog.Text, fmt.Sprint(123, "321"))

		level.logFFunction("%d 123 %s", 123, "321")
		assert.Equal(t, w.lastLog.Level, level.l, "Level should be same.")
		assert.Less(t, time.Since(w.lastLog.Time).Nanoseconds(), time.Millisecond.Nanoseconds(), "Time difference should be less than 1 ms.")
		_, _, line, _ = runtime.Caller(0)
		assert.Equal(t, w.lastLog.Caller, fmt.Sprint("github.com/leoleoasd/EduOJBackend/base/logging.TestLogger:", line-3), "Caller should be this function.")
		assert.Equal(t, w.lastLog.Text, fmt.Sprintf("%d 123 %s", 123, "321"))
	}
}
