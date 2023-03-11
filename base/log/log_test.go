package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	levels := []struct {
		Level
		string
	}{
		{
			DEBUG,
			"DEBUG",
		},
		{
			INFO,
			"INFO",
		},
		{
			WARNING,
			"WARNING",
		},
		{
			ERROR,
			"ERROR",
		},
		{
			FATAL,
			"FATAL",
		},
		{
			100,
			"",
		},
	}
	for _, level := range levels {
		assert.Equal(t, level.Level.String(), level.string)
	}
}
