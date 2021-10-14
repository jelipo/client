package main

import (
	"client/api"
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
	httpApi := api.NewRunnerHttpApi(
		config.GlobalConfig.Server.Address,
		config.GlobalConfig.Server.RunnerId,
		config.GlobalConfig.Server.Token,
	)
	return RunnerClient{
		jobManager:  work.NewWorkerManager(config.GlobalConfig.Server.MaxJobNum),
		runnerAlive: NewRunnerAlive(&httpApi),
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
}

// 向服务端请求
func (runnerManager *RunnerClient) requestServer() {
	manageStatus, workersStatus := runnerManager.jobManager.ReadStatus()
	acceptJobs, denyJobs := runnerManager.jobManager.GetAndCleanAcceptDenyRunningJobIds()
	aliveResponse, err := runnerManager.runnerAlive.alive(&manageStatus, workersStatus, acceptJobs, denyJobs)
	if err != nil {
		logrus.Warning("alive error:" + err.Error())
		return
	}
	newJobs := aliveResponse.NewJobs
	logrus.Info("Get ", len(newJobs), " jobs")
	runnerManager.handleAliveNewJobs(newJobs)
}

// 处理新的Jobs
func (runnerManager *RunnerClient) handleAliveNewJobs(newJobs []api.NewJob) {
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
