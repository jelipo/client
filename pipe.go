package main

import "time"

type PipeManager struct {
}

func run() {

	for true {
		time.Sleep(time.Duration(2000) * time.Millisecond)
	}
}
