package log

import (
	"github.com/leoleoasd/EduOJBackend/base/event"
	"github.com/leoleoasd/EduOJBackend/base/logging"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogEvent(t *testing.T) {
	lastLog := logging.Log{}
	event.RegisterListener("test_log_event", func(arg EventArgs) {
		lastLog = arg
	})
	log := logging.Log{
		Level:   logging.DEBUG,
		Time:    time.Now(),
		Message: "123",
		Caller:  "233",
	}
	event.FireEvent("test_log_event", log)
	assert.Equal(t, log, lastLog)
}
