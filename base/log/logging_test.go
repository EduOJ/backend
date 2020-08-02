package log

import (
	"errors"
	"fmt"
	"github.com/kami-zh/go-capturer"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type fakeLogger struct {
	ready              bool
	disabled           bool
	lastFunctionCalled string
}

func (f *fakeLogger) addWriter(writer _writer) {}

func (f *fakeLogger) removeWriter(t reflect.Type) {}

func (f *fakeLogger) isReady() bool {
	return f.ready
}

func (f *fakeLogger) setReady() {
	f.ready = true
}
func (f *fakeLogger) Disabled() bool {
	return f.disabled
}

func (f *fakeLogger) Disable() {
	f.disabled = true
}

func (f *fakeLogger) Debug(items ...interface{}) {
	f.lastFunctionCalled = "Debug"
}

func (f *fakeLogger) Info(items ...interface{}) {
	f.lastFunctionCalled = "Info"
}

func (f *fakeLogger) Warning(items ...interface{}) {
	f.lastFunctionCalled = "Warning"
}

func (f *fakeLogger) Error(items ...interface{}) {
	f.lastFunctionCalled = "Error"
}

func (f *fakeLogger) Fatal(items ...interface{}) {
	f.lastFunctionCalled = "Fatal"
}

func (f *fakeLogger) Debugf(format string, items ...interface{}) {
	f.lastFunctionCalled = "Debugf"
}

func (f *fakeLogger) Infof(format string, items ...interface{}) {
	f.lastFunctionCalled = "Infof"
}

func (f *fakeLogger) Warningf(format string, items ...interface{}) {
	f.lastFunctionCalled = "Warningf"
}

func (f *fakeLogger) Errorf(format string, items ...interface{}) {
	f.lastFunctionCalled = "Errorf"
}

func (f *fakeLogger) Fatalf(format string, items ...interface{}) {
	f.lastFunctionCalled = "Fatalf"
}

func TestLogFunctions(t *testing.T) {
	oldLogger := logger0
	t.Cleanup(func() {
		logger0 = oldLogger
	})
	f := &fakeLogger{}
	logger0 = f
	tests := []struct {
		function interface{}
		name     string
	}{
		{
			Debug,
			"Debug",
		},
		{
			Info,
			"Info",
		},
		{
			Warning,
			"Warning",
		},
		{
			Error,
			"Error",
		},
		{
			Fatal,
			"Fatal",
		},
		{
			Debugf,
			"Debugf",
		},
		{
			Infof,
			"Infof",
		},
		{
			Warningf,
			"Warningf",
		},
		{
			Errorf,
			"Errorf",
		},
		{
			Fatalf,
			"Fatalf",
		},
	}
	for _, test := range tests {
		t.Run("testLogFunction"+test.name, func(t *testing.T) {
			if _, ok := test.function.(func(...interface{})); ok {
				test.function.(func(...interface{}))()
			} else {
				test.function.(func(string, ...interface{}))("")
			}
			assert.Equal(t, test.name, f.lastFunctionCalled)
		})
	}
}

func TestInitFromConfigFail(t *testing.T) {
	oldLogger := logger0
	t.Cleanup(func() {
		logger0 = oldLogger
	})
	tests := []struct {
		config.Node
		error
	}{
		{
			&config.MapNode{}, errors.New("log configuration should be an array"),
		},
		{
			nil, errors.New("log configuration should be an array"),
		},
		{
			&config.SliceNode{S: []config.Node{
				&config.MapNode{M: map[string]config.Node{
					"name": config.StringNode("invalid_writer_name"),
				}},
			}}, errors.New("invalid writer name"),
		},
		{
			&config.SliceNode{S: []config.Node{
				&config.MapNode{M: map[string]config.Node{
					"name": config.IntNode(123),
				}},
			}}, errors.New("invalid writer name"),
		},
		{
			&config.SliceNode{S: []config.Node{
				&config.MapNode{M: map[string]config.Node{}},
			}}, errors.New("writer configuration should contain name"),
		},
		{
			&config.SliceNode{S: []config.Node{
				config.IntNode(123),
			}}, errors.New("writer configuration should be a map"),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint("testInit_", i), func(t *testing.T) {
			l := &logger{}
			logger0 = l
			err := InitFromConfig(test.Node)
			if test.error != nil && err != nil {
				assert.Equal(t, test.error.Error(), err.Error())
			} else {
				assert.Equal(t, test.error, err)
			}
		})
	}
}

func TestInitFromConfigSuccess(t *testing.T) {
	t.Cleanup(database.SetupDatabaseForTest())
	t.Cleanup(exit.SetupExitForTest())
	oldLogger := logger0
	t.Cleanup(func() {
		logger0 = oldLogger
	})
	l := &logger{}
	logger0 = l
	assert.Equal(t, false, l.ready)
	err := InitFromConfig(
		&config.SliceNode{S: []config.Node{
			&config.MapNode{M: map[string]config.Node{
				"name":  config.StringNode("console"),
				"level": config.StringNode("InFO"),
			}},
			&config.MapNode{M: map[string]config.Node{
				"name":  config.StringNode("database"),
				"level": config.StringNode("InFO"),
			}},
			&config.MapNode{M: map[string]config.Node{
				"name":  config.StringNode("event"),
				"level": config.StringNode("InFO"),
			}},
		}})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, l.ready)
	err = InitFromConfig(&config.SliceNode{S: []config.Node{}})
	assert.EqualError(t, err, "already initialized")
	exit.Close()
	exit.QuitWG.Wait()
}

func TestLogging_Disable(t *testing.T) {
	oldLogger := logger0
	t.Cleanup(func() {
		logger0 = oldLogger
	})
	l := &logger{}
	logger0 = l
	assert.Equal(t, false, l.disabled)

	output := capturer.CaptureOutput(func() {
		l.Debug("test")
	})
	assert.NotEqual(t, "", output)

	Disable()
	assert.Equal(t, true, Disabled())
	output = capturer.CaptureOutput(func() {
		l.Debug("test")
	})
	assert.Equal(t, "", output)
}
