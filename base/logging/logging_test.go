package logging

import (
	"errors"
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type fakeLogger struct {
	ready              bool
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
		Debug = logger0.Debug
		Info = logger0.Info
		Warning = logger0.Warning
		Error = logger0.Error
		Fatal = logger0.Fatal
		Debugf = logger0.Debugf
		Infof = logger0.Infof
		Warningf = logger0.Warningf
		Errorf = logger0.Errorf
		Fatalf = logger0.Fatalf
	})
	f := &fakeLogger{}
	logger0 = f
	Debug = logger0.Debug
	Info = logger0.Info
	Warning = logger0.Warning
	Error = logger0.Error
	Fatal = logger0.Fatal
	Debugf = logger0.Debugf
	Infof = logger0.Infof
	Warningf = logger0.Warningf
	Errorf = logger0.Errorf
	Fatalf = logger0.Fatalf
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
		Debug = logger0.Debug
		Info = logger0.Info
		Warning = logger0.Warning
		Error = logger0.Error
		Fatal = logger0.Fatal
		Debugf = logger0.Debugf
		Infof = logger0.Infof
		Warningf = logger0.Warningf
		Errorf = logger0.Errorf
		Fatalf = logger0.Fatalf
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
			Debug = logger0.Debug
			Info = logger0.Info
			Warning = logger0.Warning
			Error = logger0.Error
			Fatal = logger0.Fatal
			Debugf = logger0.Debugf
			Infof = logger0.Infof
			Warningf = logger0.Warningf
			Errorf = logger0.Errorf
			Fatalf = logger0.Fatalf
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
	oldLogger := logger0
	t.Cleanup(func() {
		logger0 = oldLogger
		Debug = logger0.Debug
		Info = logger0.Info
		Warning = logger0.Warning
		Error = logger0.Error
		Fatal = logger0.Fatal
		Debugf = logger0.Debugf
		Infof = logger0.Infof
		Warningf = logger0.Warningf
		Errorf = logger0.Errorf
		Fatalf = logger0.Fatalf
	})
	l := &logger{}
	logger0 = l
	Debug = logger0.Debug
	Info = logger0.Info
	Warning = logger0.Warning
	Error = logger0.Error
	Fatal = logger0.Fatal
	Debugf = logger0.Debugf
	Infof = logger0.Infof
	Warningf = logger0.Warningf
	Errorf = logger0.Errorf
	Fatalf = logger0.Fatalf
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
}
