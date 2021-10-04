package main

import (
	"client/config"
	"client/work"
	"github.com/sirupsen/logrus"
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
		logrus.Info("Start get new jobs from server")
		runnerManager.requestServer()
		var sleepMills = 2000
		time.Sleep(time.Duration(sleepMills) * time.Millisecond)
	}

	//logrus.Info("Start get new jobs from server")
	//runnerManager.requestServer()
	//var sleepMills = 2000
	//time.Sleep(time.Duration(sleepMills) * time.Second)

}

func (runnerManager *RunnerClient) requestServer() {
	manageStatus, workersStatus := runnerManager.jobManager.ReadStatus()
	acceptJobs, denyJobs := runnerManager.jobManager.GetAndCleanAcceptDenyRunningJobIds()
	aliveResponse, err := runnerManager.runnerAlive.alive(&manageStatus, workersStatus, acceptJobs, denyJobs)
	if err != nil {
		logrus.Info("alive error:" + err.Error())
		return
	}
	newJobs := aliveResponse.NewJobs
	logrus.Info("Get ", len(newJobs), " jobs")
	runnerManager.handleAliveNewJobs(newJobs)
}

// 处理新的Jobs
func (runnerManager *RunnerClient) handleAliveNewJobs(newJobs []work.NewJob) {
	var acceptRunningJobIds []string
	var denyRunningJobIds []string
	for _, job := range newJobs {
		err := runnerManager.jobManager.AddNewJob(job.JobRunningId, job.Sources, &job)
		if err != nil {
			logrus.Info("add new job error" + err.Error())
			denyRunningJobIds = append(denyRunningJobIds, job.JobRunningId)
		} else {
			logrus.Info("add new job success, JobRunningId:" + job.JobRunningId)
			acceptRunningJobIds = append(acceptRunningJobIds, job.JobRunningId)
		}
	}
	runnerManager.jobManager.AddAcceptDenyRunningJobId(acceptRunningJobIds, denyRunningJobIds)
}
