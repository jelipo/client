package util

import "github.com/sirupsen/logrus"

type AsyncRunFlag struct {
	channel   chan int
	haveError bool
}

func NewAsyncRunFlag(fn func() error) AsyncRunFlag {
	flag := AsyncRunFlag{
		channel:   make(chan int, 1),
		haveError: false,
	}
	go flag.run(fn)
	return flag
}

func (flag *AsyncRunFlag) IsDone() bool {
	return len(flag.channel) > 0
}

func (flag *AsyncRunFlag) HaveError() bool {
	return flag.haveError
}

func (flag *AsyncRunFlag) run(fn func() error) {
	err := fn()
	if err != nil {
		flag.haveError = true
		logrus.Info("RunningJob happened a error:" + err.Error())
		return
	} else {
		flag.haveError = false
	}
	flag.channel <- 1
}
