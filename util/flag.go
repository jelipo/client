package util

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type AsyncRunFlag struct {
	rwMutex     sync.RWMutex
	doneChannel bool
	errChannel  bool
}

func NewAsyncRunFlag(fn func() error) *AsyncRunFlag {
	flag := AsyncRunFlag{
		rwMutex:     sync.RWMutex{},
		doneChannel: false,
		errChannel:  false,
	}
	go flag.run(fn)
	return &flag
}

func (flag *AsyncRunFlag) IsDone() bool {
	flag.rwMutex.RLock()
	defer flag.rwMutex.RUnlock()
	return flag.doneChannel
}

func (flag *AsyncRunFlag) HaveError() bool {
	flag.rwMutex.RLock()
	defer flag.rwMutex.RUnlock()
	return flag.errChannel
}

func (flag *AsyncRunFlag) run(fn func() error) {
	err := fn()
	flag.rwMutex.Lock()
	defer flag.rwMutex.Unlock()
	if err != nil {
		flag.errChannel = true
		logrus.Info("RunningJob happened a error:" + err.Error())
	}
	flag.doneChannel = true
}
