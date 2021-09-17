package main

import (
	"client/config"
	"client/work"
	"log"
	"time"
)

type RunnerClient struct {
	jobManager  work.JobManager
	runnerAlive RunnerAlive
}

func NewRunnerClient() RunnerClient {
	return RunnerClient{
		jobManager:  work.NewWorkerManager(config.GlobalConfig.Server.MaxWorkerNum),
		runnerAlive: NewRunnerAlive(),
	}
}

func (runnerManager *RunnerClient) run() {
	// 单线程循环向服务端请求
	for true {
		log.Println("Start get new jobs from server")
		runnerManager.requestServer()
		var sleepMills = 2000
		time.Sleep(time.Duration(sleepMills) * time.Millisecond)
	}

	//log.Println("Start get new jobs from server")
	//runnerManager.requestServer()
	//var sleepMills = 2000
	//time.Sleep(time.Duration(sleepMills) * time.Second)

}

func (runnerManager *RunnerClient) requestServer() {
	manageStatus, workersStatus := runnerManager.jobManager.ReadStatus()
	acceptJobs, denyJobs := runnerManager.jobManager.GetAndCleanAcceptDenyRunningJobIds()
	aliveResponse, err := runnerManager.runnerAlive.alive(&manageStatus, workersStatus, acceptJobs, denyJobs)
	if err != nil {
		log.Println("alive error:" + err.Error())
		return
	}
	newJobs := aliveResponse.NewJobs
	log.Println("Get ", len(newJobs), " jobs")
	runnerManager.handleAliveNewJobs(newJobs)
}

// 处理新的Jobs
func (runnerManager *RunnerClient) handleAliveNewJobs(newJobs []work.NewJob) {
	var acceptRunningJobIds []string
	var denyRunningJobIds []string
	for _, job := range newJobs {
		err := runnerManager.jobManager.AddNewJob(job.JobRunningId, job.Sources, &job)
		if err != nil {
			log.Println("add new job error" + err.Error())
			denyRunningJobIds = append(denyRunningJobIds, job.JobRunningId)
		} else {
			log.Println("add new job success, JobRunningId:" + job.JobRunningId)
			acceptRunningJobIds = append(acceptRunningJobIds, job.JobRunningId)
		}
	}
	runnerManager.jobManager.AddAcceptDenyRunningJobId(acceptRunningJobIds, denyRunningJobIds)
}
