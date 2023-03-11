package log

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
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
	l.setReady()
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
		assert.Equal(t, w.lastLog.Message, fmt.Sprint(123, "321"))

		level.logFFunction("%d 123 %s", 123, "321")
		assert.Equal(t, w.lastLog.Level, level.l, "Level should be same.")
		assert.Equal(t, w.lastLog.Message, fmt.Sprintf("%d 123 %s", 123, "321"))
	}
}

func TestRemoveLogger(t *testing.T) {
	var l _logger
	l = &logger{}
	l.setReady()
	w1 := &fakeWriter{}
	w2 := &fakeWriter{}
	w3 := &consoleWriter{}
	l.addWriter(w1)
	l.addWriter(w2)
	l.addWriter(w3)
	assert.Equal(t, w1, l.(*logger).writers[0], "Writer should be fake writer.")
	assert.Equal(t, w2, l.(*logger).writers[1], "Writer should be fake writer.")
	assert.Equal(t, w3, l.(*logger).writers[2], "Writer should be fake writer.")
	l.removeWriter(reflect.TypeOf((*fakeWriter)(nil)))
	assert.Equal(t, l.(*logger).writers, []_writer{w3}, "Should not have any writers here.")
}

func TestLogWithLevelString(t *testing.T) {
	var l logger
	w := &fakeWriter{}
	l.addWriter(w)
	l.setReady()
	levels := []Level{
		DEBUG,
		INFO,
		WARNING,
		ERROR,
		FATAL,
	}
	for _, level := range levels {
		l.logWithLevelString(level, "test")
		assert.Equal(t, w.lastLog.Level, level)
		assert.Equal(t, w.lastLog.Message, "test")
		assert.Regexp(t, "^testing\\.tRunner:[0-9]+$", w.lastLog.Caller)
	}
}

func TestLogNotReady(t *testing.T) {
	var l _logger
	l = &logger{}

	levels := []struct {
		Level
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
		output := capturer.CaptureOutput(func() {
			level.logFunction("sample log output")
		})
		txt := fmt.Sprintf("[%s] ▶ %s %s\n",
			"github.com/kami-zh/go-capturer.(*Capturer).capture:55",
			level.Level.String(),
			"sample log output")
		assert.Equal(t, txt, output[10:])
		output = capturer.CaptureOutput(func() {
			level.logFFunction("sample log output")
		})
		txt = fmt.Sprintf("[%s] ▶ %s %s\n",
			"github.com/kami-zh/go-capturer.(*Capturer).capture:55",
			level.Level.String(),
			"sample log output")
		assert.Equal(t, txt, output[10:])
	}
}

func TestIsReady(t *testing.T) {
	var l _logger
	l = &logger{}
	assert.Equal(t, false, l.isReady())
	l.setReady()
	assert.Equal(t, true, l.isReady())
}

func TestLogger_Disable(t *testing.T) {
	var l = &logger{}
	assert.Equal(t, false, l.disabled)

	output := capturer.CaptureOutput(func() {
		l.Debug("test")
	})
	assert.NotEqual(t, "", output)

	l.Disable()
	assert.Equal(t, true, l.Disabled())
	output = capturer.CaptureOutput(func() {
		l.Debug("test")
	})
	assert.Equal(t, "", output)
}
