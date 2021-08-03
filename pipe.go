package main

import (
	"client/config"
	"client/work"
	"fmt"
	"time"
)

type PipeManager struct {
	workerManager work.Manager
	runnerAlive   RunnerAlive
}

func NewPipeManager() PipeManager {
	return PipeManager{
		workerManager: work.NewWorkerManager(config.GlobalConfig.Server.MaxWorkerNum),
		runnerAlive:   NewRunnerAlive(),
	}
}

func (manager *PipeManager) run() {
	for true {
		fmt.Println("Start get new job from server")
		manager.aliveServer()
		time.Sleep(time.Duration(2000) * time.Millisecond)
	}
}

func (manager *PipeManager) aliveServer() {
	manageStatus, workersStatus := manager.workerManager.ReadStatus()
	alive, err := manager.runnerAlive.alive(&manageStatus, workersStatus, make([]string, 0))
	if err != nil {
		fmt.Println("alive error:" + err.Error())
		return
	}
	newJobs := alive.NewJobs
	fmt.Println("Get ", len(newJobs), " jobs")
	for _, job := range newJobs {
		err := manager.workerManager.AddNewJob(job.JobRunningId, job.Sources, &job.NewWork)
		if err != nil {
			// TODO
			fmt.Println("add new job error" + err.Error())
			return
		}
	}
}
