package util

import (
	"sync"
)

type ThreadGroup struct {
	workGroup sync.WaitGroup
	lock      bool
}

func NewGoPackage() ThreadGroup {
	return ThreadGroup{
		workGroup: sync.WaitGroup{},
		lock:      false,
	}
}

func (pack *ThreadGroup) AddAndRun(a func()) {
	if pack.lock {
		return
	}
	go pack.run(a)
	pack.workGroup.Add(1)
}

func (pack *ThreadGroup) WaitAllDone() {
	pack.lock = true
	pack.workGroup.Wait()
}

func (pack *ThreadGroup) run(a func()) {
	a()
	pack.workGroup.Done()
}
