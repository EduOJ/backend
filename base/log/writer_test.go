package log

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kami-zh/go-capturer"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/event"
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)
import _ "github.com/jinzhu/gorm/dialects/sqlite"

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

func TestEventWriter(t *testing.T) {
	lastLog := Log{}
	done := make(chan struct{})
	event.RegisterListener("log", func(arg Log) {
		lastLog = arg
		done <- struct{}{}
	})
	w := &eventWriter{}
	log := Log{
		Level:   DEBUG,
		Time:    time.Now(),
		Message: "123",
		Caller:  "233",
	}
	w.log(log)
	<-done
	assert.Equal(t, log, lastLog)
}

func TestDatabaseWriter(t *testing.T) {
	oldDB := base.DB
	base.DB, _ = gorm.Open("sqlite3", ":memory:")
	database.Migrate()
	oldWG := exit.QuitWG
	exit.QuitWG = sync.WaitGroup{}
	oldClose := exit.Close
	oldContext := exit.BaseContext
	exit.BaseContext, exit.Close = context.WithCancel(context.Background())
	t.Cleanup(func() {
		base.DB = oldDB
		exit.QuitWG = oldWG
		exit.Close = oldClose
		exit.BaseContext = oldContext
	})
	log := Log{
		Level:   DEBUG,
		Time:    time.Now(),
		Message: "123",
		Caller:  "233",
	}
	w := databaseWriter{}
	w.init()
	w.log(log)
	exit.Close()
	exit.QuitWG.Wait()
	lm := models.Log{}
	base.DB.First(&lm)
	ll := int(DEBUG)
	assert.Equal(t, models.Log{
		ID:        lm.ID,
		Level:     &ll,
		Message:   "123",
		Caller:    "233",
		CreatedAt: lm.CreatedAt,
	}, lm)
}
