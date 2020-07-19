package exit

import (
	"context"
	"sync"
)

var BaseContext, Close = context.WithCancel(context.Background())
var QuitWG = sync.WaitGroup{}

var testExitLock = sync.Mutex{}

//noinspection GoVetCopyLock
func SetupExitForTest() func() {
	testExitLock.Lock()
	oldWG := QuitWG
	oldClose := Close
	oldContext := BaseContext
	QuitWG = sync.WaitGroup{}
	BaseContext, Close = context.WithCancel(context.Background())
	return func() {
		QuitWG = oldWG
		Close = oldClose
		BaseContext = oldContext
		testExitLock.Unlock()
	}
}
