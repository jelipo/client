package main

import (
	"client/config"
	"client/work"
)

type RunnerAlive struct {
	serverAddress string
	token         string
	//workStatusCache map[string]WorkStatusCache
}

type WorkStatusCache struct {
}

func NewRunnerAlive(serverAddr string, token string) RunnerAlive {
	return RunnerAlive{
		serverAddress: serverAddr,
		token:         token,
	}
}

func (alive *RunnerAlive) alive(
	runnerStatus *work.ManagerStatus,
	workersStatus map[string]work.WorkerOutStatus,
	acceptJobs []string,
) AliveResponse {

	aliveRequest := AliveRequest{
		HostStatus:   HostStatus{},
		RunnerStatus: *runnerStatus,
		JobsStatus:   changeStatus(workersStatus),
		AcceptJobs:   acceptJobs,
	}
	config.GlobalConfig.Server.
	// TODO http
	return AliveResponse{}
}

func changeStatus(workerStatus map[string]work.WorkerOutStatus) []JobsStatus {
	var jobsStatus = make([]JobsStatus, len(workerStatus))
	for _, workStatus := range workerStatus {
		jobsStatus = append(jobsStatus, JobsStatus{
			JobRunningId: workStatus.JobRunningId,
			AtomLogs:     workStatus.AtomLogs,
			Done:         workStatus.Done,
		})
	}
	return jobsStatus
}

type NewJob struct {
	Sources []work.Source `json:"sources"`
	NewWork work.NewWork  `json:"newWork"`
}

type AliveResponse struct {
	NewJobs []NewJob `json:"newJobs"`
}

type AliveRequest struct {
	HostStatus   HostStatus         `json:"hostStatus"`
	RunnerStatus work.ManagerStatus `json:"runnerStatus"`
	JobsStatus   []JobsStatus       `json:"jobsStatus"`
	AcceptJobs   []string           `json:"acceptJobs"`
}

type HostStatus struct {
	// TODO CPU/Memory/Disk info
}

type JobsStatus struct {
	JobRunningId string         `json:"jobRunningId"`
	AtomLogs     []work.AtomLog `json:"atomLogs"`
	Done         bool           `json:"done"`
}
