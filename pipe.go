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

func (manager *PipeManager) run() {
	for true {
		manageStatus, workersStatus := manager.workerManager.ReadStatus()

		time.Sleep(time.Duration(2000) * time.Millisecond)
	}
}
