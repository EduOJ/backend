package base

import (
	"time"
)

// Log levels.
type Level int

const (
	DEBUG   Level = iota // Debug information.
	INFO                 // Running information.
	WARNING              // Warnings: notable, but the process wont fail because of this.
	ERROR                // Errors: a process (request) fails because of this.
	FATAL                // Fatal: multiple process(request) fails because of this.
)

// Inner log interface.
type _log interface {
	Level() Level
	Time() time.Time
	Text() string
	Function() string
}

// Log
type Log struct {
	Level `json:"level"`
	Time time.Time `json:"time"`
	Text string `json:"text"`
	Function string `json:"function"`
}

// Logger should only be used in this package.
// Other packages should use hooks to receive logs.
type _logger interface {

}





func (l Level) String() string {
	switch l{
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default :
		return "FATAL"
	}
}
