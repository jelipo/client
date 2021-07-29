package util

type AsyncCollect struct {
	//workGroup sync.WaitGroup
	asyncChan chan bool
	num       int
	lock      bool
}

func NewAsyncCollect() AsyncCollect {
	return AsyncCollect{
		asyncChan: make(chan bool, 1024),
		num:       0,
		lock:      false,
	}
}

func (pack *AsyncCollect) AddAndRun(a func()) {
	if pack.lock {
		return
	}
	go pack.run(a)
	pack.num = pack.num + 1
	pack.asyncChan <- true
}

func (pack *AsyncCollect) IsAllDone() bool {
	return len(pack.asyncChan) == 0
}

func (pack *AsyncCollect) run(a func()) {
	a()
	_ = <-pack.asyncChan
}
