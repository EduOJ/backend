package logging

import (
	"time"
)

// Log levels.
type Level int

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		// Shouldn't reach here.
		return ""
	}
}

const (
	DEBUG   Level = iota // Debug information.
	INFO                 // Running information.
	WARNING              // Warnings: notable, but the process wont fail because of this.
	ERROR                // Errors: a process (request) fails because of this.
	FATAL                // Fatal: multiple process(request) fails because of this.
)

// The log struct
type Log struct {
	Level  Level     `json:"level"`  // The level of this log.
	Time   time.Time `json:"time"`   // The time of this log.
	Text   string    `json:"text"`   // The text of this log.
	Caller string    `json:"caller"` // The function produces this log.
}
