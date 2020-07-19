package exit

import (
	"context"
	"sync"
)

var BaseContext, Close = context.WithCancel(context.Background())
var QuitWG = sync.WaitGroup{}

//noinspection GoVetCopyLock
func SetupExitForTest() func() {
	oldWG := QuitWG
	QuitWG = sync.WaitGroup{}
	oldClose := Close
	oldContext := BaseContext
	BaseContext, Close = context.WithCancel(context.Background())
	return func() {
		QuitWG = oldWG
		Close = oldClose
		BaseContext = oldContext
	}
}
