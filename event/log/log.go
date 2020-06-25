package log

import "github.com/leoleoasd/EduOJBackend/base/logging"

// EventArgs is the arguments of "log" event.
// Only contains the log itself.
type EventArgs = logging.Log

// EventRst is the result of "log" event.
// Cause there is no need for result, this
// is a empty struct.
type EventRst struct{}
