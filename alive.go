package main

import (
	"client/api"
	"client/work"
	"io/ioutil"
	"net/http"
)

type RunnerAlive struct {
	runnerHttpApi *api.RunnerHttpApi
	//workStatusCache map[string]WorkStatusCache
}

type WorkStatusCache struct {
}

func NewRunnerAlive(runnerHttpApi *api.RunnerHttpApi) RunnerAlive {
	return RunnerAlive{
		runnerHttpApi: runnerHttpApi,
	}
}

func (alive *RunnerAlive) alive(
	runnerStatus *work.ManagerStatus,
	workersStatus map[string]work.JobOutStatus,
	acceptRunningJobIds []string,
	denyRunningJobIds []string,
) (*api.AliveResponse, error) {
	aliveRequest := api.AliveRequest{
		HostStatus:   api.HostStatus{},
		RunnerStatus: api.RunnerStatus{RunningNum: runnerStatus.RunningNum},
		JobsStatus:   changeJobStatusToRequest(workersStatus),
		AcceptJobs:   acceptRunningJobIds,
		DenyJobs:     denyRunningJobIds,
	}
	aliveResponse, err := alive.runnerHttpApi.AliveToServer(&aliveRequest)
	if err != nil {
		return nil, err
	}
	return aliveResponse, nil
}

func changeJobStatusToRequest(workerStatus map[string]work.JobOutStatus) []api.JobsStatus {
	var jobsStatus []api.JobsStatus
	for _, workStatus := range workerStatus {
		status := api.JobsStatus{
			JobRunningId: workStatus.JobRunningId,
			AtomLogs:     workStatus.AtomLogs,
			Finished:     workStatus.Finished,
		}
		if workStatus.Success {
			status.FinishedStatus = "SUCCESS"
		}
		jobsStatus = append(jobsStatus, status)
	}
	return jobsStatus
}

func readBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseBodyBytes, nil
}
