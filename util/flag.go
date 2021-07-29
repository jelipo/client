package util

type AsyncRunFlag struct {
	channel chan int
}

func NewAsyncRunFlag(fn func()) AsyncRunFlag {
	flag := AsyncRunFlag{
		channel: make(chan int, 1),
	}
	go flag.run(fn)
	return flag
}

func (flag *AsyncRunFlag) IsDone() bool {
	return len(flag.channel) > 0
}

func (flag *AsyncRunFlag) run(fn func()) {
	fn()
	flag.channel <- 1
}
