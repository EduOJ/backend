package log

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/event"
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/database/models"
)

// Writers should only be used in this package.
// Other packages should use listeners to receive logs.
type _writer interface {
	log(log Log)
}

// Writes to the console.
type consoleWriter struct {
	Level
}

// Writes to the database for reading from web.
type databaseWriter struct {
	Level
	queue chan Log
}

// Calling log listeners.
type eventWriter struct{}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

var (
	colors = []string{
		FATAL:   colorSeq(colorMagenta),
		ERROR:   colorSeq(colorRed),
		WARNING: colorSeq(colorYellow),
		INFO:    colorSeq(colorGreen),
		DEBUG:   colorSeq(colorCyan),
	}
)

func colorSeq(color int) string {
	return fmt.Sprintf("\033[%dm", color)
}

func (w *consoleWriter) log(l Log) {
	if l.Level >= w.Level {
		fmt.Printf("%s[%s][%s] ▶ %s\u001B[0m %s\n",
			colors[l.Level],
			l.Time.Format("15:04:05"),
			l.Caller,
			l.Level.String(),
			l.Message)
	}
}

func (w *databaseWriter) log(l Log) {
	// avoid blocking the main thread.
	if l.Level >= w.Level {
		select {
		case w.queue <- l:
		default:
		}
	}
}

func (w *databaseWriter) init() {
	w.queue = make(chan Log, 100)
	exit.QuitWG.Add(1)
	go func() {
		for {
			select {
			case l := <-w.queue:
				lm := models.Log{
					Model: gorm.Model{
						CreatedAt: l.Time,
						UpdatedAt: l.Time,
					},
					Level:   int(l.Level),
					Message: l.Message,
					Caller:  l.Caller,
				}
				base.DB.Create(&lm)
			case <-exit.BaseContext.Done():
				select {
				case l := <-w.queue:
					lm := models.Log{
						Model: gorm.Model{
							CreatedAt: l.Time,
							UpdatedAt: l.Time,
						},
						Level:   int(l.Level),
						Message: l.Message,
						Caller:  l.Caller,
					}
					base.DB.Create(&lm)
				default:
					exit.QuitWG.Done()
					return
				}
			}
		}
	}()
}

func (w *eventWriter) log(l Log) {
	go event.FireEvent("log", l)
}
