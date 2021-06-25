package main

import (
	"sync"
)

type GoPackage struct {
	workGroup sync.WaitGroup
	lock      bool
}

func NewGoPackage() GoPackage {
	return GoPackage{
		workGroup: sync.WaitGroup{},
		lock:      false,
	}
}

func (pack *GoPackage) AddAndRun(a func()) {
	if pack.lock {
		return
	}
	go pack.run(a)
	pack.workGroup.Add(1)
}

func (pack *GoPackage) WaitAllDone() {
	pack.lock = true
	pack.workGroup.Wait()
}

func (pack *GoPackage) run(a func()) {
	a()
	pack.workGroup.Done()
}
