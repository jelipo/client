package util

type GoroutinesFlag struct {
	channel chan int
}

func NewGoroutinesFlagAndRun(fn func()) GoroutinesFlag {
	flag := GoroutinesFlag{
		channel: make(chan int, 1),
	}
	go flag.run(fn)
	return flag
}

func (flag *GoroutinesFlag) IsDone() bool {
	return len(flag.channel) > 0
}

func (flag *GoroutinesFlag) run(fn func()) {
	fn()
	flag.channel <- 1
}
