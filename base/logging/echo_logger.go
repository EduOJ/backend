package logging

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"io"
	. "os"
)

// EchoLogger is a fake logger for echo.
type EchoLogger struct {
	prefix string
}

func (e *EchoLogger) Output() io.Writer {
	return Stdout
}

func (e *EchoLogger) SetOutput(w io.Writer) {
	// do nothing
}

func (e *EchoLogger) Prefix() string {
	return e.prefix
}

func (e *EchoLogger) SetPrefix(p string) {
	e.prefix = p
}

func (e *EchoLogger) Level() log.Lvl {
	return 0
}

func (e *EchoLogger) SetLevel(v log.Lvl) {
	// do nothing
}

func (e *EchoLogger) SetHeader(h string) {
	// do nothing
}

func (e *EchoLogger) Print(i ...interface{}) {
	Info(i...)
}

func (e *EchoLogger) Printf(format string, args ...interface{}) {
	Infof(format, args...)
}

func (e *EchoLogger) Printj(j log.JSON) {
	b, _ := json.Marshal(j)
	Info(string(b))
}

func (e *EchoLogger) Debug(i ...interface{}) {
	Debug(i...)
}

func (e *EchoLogger) Debugf(format string, args ...interface{}) {
	Debugf(format, args...)
}

func (e *EchoLogger) Debugj(j log.JSON) {
	b, _ := json.Marshal(j)
	Debug(string(b))
}

func (e *EchoLogger) Info(i ...interface{}) {
	Info(i...)
}

func (e *EchoLogger) Infof(format string, args ...interface{}) {
	Infof(format, args...)
}

func (e *EchoLogger) Infoj(j log.JSON) {
	b, _ := json.Marshal(j)
	Info(string(b))
}

func (e *EchoLogger) Warn(i ...interface{}) {
	Warning(i...)
}

func (e *EchoLogger) Warnf(format string, args ...interface{}) {
	Warningf(format, args...)
}

func (e *EchoLogger) Warnj(j log.JSON) {
	b, _ := json.Marshal(j)
	Warning(string(b))
}

func (e *EchoLogger) Error(i ...interface{}) {
	Error(i...)
}

func (e *EchoLogger) Errorf(format string, args ...interface{}) {
	Errorf(format, args...)
}

func (e *EchoLogger) Errorj(j log.JSON) {
	b, _ := json.Marshal(j)
	Error(string(b))
}

func (e *EchoLogger) Fatal(i ...interface{}) {
	Fatal(i...)
}

func (e *EchoLogger) Fatalf(format string, args ...interface{}) {
	Fatalf(format, args...)
}

func (e *EchoLogger) Fatalj(j log.JSON) {
	b, _ := json.Marshal(j)
	Fatal(string(b))
}

func (e *EchoLogger) Panic(i ...interface{}) {
	Fatal(i...)
}

func (e *EchoLogger) Panicf(format string, args ...interface{}) {
	Fatalf(format, args...)
}

func (e *EchoLogger) Panicj(j log.JSON) {
	b, _ := json.Marshal(j)
	Fatal(string(b))
}
