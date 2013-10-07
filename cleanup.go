// cleanup project cleanup.go
package cleanup

import (
	"container/list"
	"errors"
	"sync"
)

type step func() error
type steps list.List

var (
	shutdownChannel = make(chan bool)
	shutdownOnce    sync.Once
)

var (
	First = (*steps)(list.New())
	Mid   = (*steps)(list.New())
	Last  = (*steps)(list.New())
)

var ShuttingDown = false

func ShutdownChannel() <-chan bool {
	return shutdownChannel
}

func Shutdown() {
	shutdownOnce.Do(func() {
		ShuttingDown = true
		shutdownChannel <- true
	})
}

func (s *steps) Append(f step) {
	(*list.List)(s).PushBack(f)
}

func (s *steps) Prepend(f step) {
	(*list.List)(s).PushFront(f)
}

func Exec() error {
	errStr := ""
	for _, l := range [3]*steps{First, Mid, Last} {
		for e := (*list.List)(l).Front(); e != nil; e = e.Next() {
			if err := e.Value.(step)(); err != nil {
				errStr += err.Error() + "\n"
			}
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	}

	return nil
}
