package logging

import (
	"fmt"
	"github.com/kami-zh/go-capturer"
	"github.com/leoleoasd/EduOJBackend/base/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConsoleWriter(t *testing.T) {
	levels := []struct {
		Level
	}{
		{
			DEBUG,
		},
		{
			INFO,
		},
		{
			WARNING,
		},
		{
			ERROR,
		},
		{
			FATAL,
		},
	}
	w := &consoleWriter{}
	ti := time.Now()
	for _, level := range levels {
		out := capturer.CaptureOutput(func() {
			w.log(Log{
				Level:   level.Level,
				Time:    ti,
				Message: "sample log output",
				Caller:  "main.main.func",
			})
		})
		txt := fmt.Sprintf("%s[%s][%s] ▶ %s\u001B[0m %s\n",
			colors[level.Level],
			ti.Format("15:04:05"),
			"main.main.func",
			level.Level.String(),
			"sample log output")
		assert.Equal(t, txt, out)
	}
}

func TestEventWriter(t *testing.T){
	lastLog := Log{}
	event.RegisterListener("log", func(arg Log) {
		lastLog = arg
	})
	w := &eventWriter{}
	log := Log{
		Level:   DEBUG,
		Time:    time.Now(),
		Message: "123",
		Caller:  "233",
	}
	w.log(log)
	assert.Equal(t, log, lastLog)
}