package exit

import (
	"context"
	"sync"
)

var BaseContext, Close = context.WithCancel(context.Background())
var QuitWG = sync.WaitGroup{}
