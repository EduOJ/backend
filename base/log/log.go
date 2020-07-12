package log

import (
	"time"
)

// Level for logs.
type Level int

/*
Debug:   debug information.
Info:    Running information.
Warning: notable, but the process wont fail because of this.
Error:   a process (request) fails because of this.
Fatal:   multiple process(request) fails because of this.
*/
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

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

var stringToLevel = map[string]Level{
	"DEBUG":   DEBUG,
	"INFO":    INFO,
	"WARNING": WARNING,
	"ERROR":   ERROR,
	"FATAL":   FATAL,
}

// Log struct contains essential information of a log.
type Log struct {
	Level   Level     `json:"level"`   // The level of this log.
	Time    time.Time `json:"time"`    // The time of this log.
	Message string    `json:"message"` // The message of this log.
	Caller  string    `json:"caller"`  // The function produces this log.
}
