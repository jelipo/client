package main

import (
	"client/config"
	"client/work"
	"fmt"
	"log"
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
		log.Println("Start get new job from server")
		manager.aliveServer()
		var sleepMills = 2000
		time.Sleep(time.Duration(sleepMills) * time.Millisecond)
	}
}

func (manager *PipeManager) aliveServer() {
	manageStatus, workersStatus := manager.workerManager.ReadStatus()
	alive, err := manager.runnerAlive.alive(&manageStatus, workersStatus, make([]string, 0))
	if err != nil {
		log.Println("alive error:" + err.Error())
		return
	}
	newJobs := alive.NewJobs
	fmt.Println("Get ", len(newJobs), " jobs")
	for _, job := range newJobs {
		err := manager.workerManager.AddNewJob(job.JobRunningId, job.Sources, &job)
		if err != nil {
			// TODO
			log.Println("add new job error" + err.Error())
			return
		}
	}
}
