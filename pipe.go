package main

import (
	"client/config"
	"client/work"
	"time"
)

type PipeManager struct {
	workerManager work.Manager
}

func NewPipeManager() PipeManager {
	return PipeManager{
		workerManager: work.NewWorkerManager(config.GlobalConfig.Server.MaxWorkerNum),
	}
}

func run() {
	for true {

		time.Sleep(time.Duration(2000) * time.Millisecond)
	}
}
