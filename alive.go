package main

import (
	"bytes"
	"client/config"
	"client/work"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RunnerAlive struct {
	httpclient http.Client
	//workStatusCache map[string]WorkStatusCache
}

type WorkStatusCache struct {
}

func NewRunnerAlive() RunnerAlive {
	return RunnerAlive{
		httpclient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (alive *RunnerAlive) alive(
	runnerStatus *work.ManagerStatus,
	workersStatus map[string]work.WorkerOutStatus,
	acceptJobs []string,
) (*AliveResponse, error) {
	aliveRequest := AliveRequest{
		HostStatus:   HostStatus{},
		RunnerStatus: *runnerStatus,
		JobsStatus:   changeStatus(workersStatus),
		AcceptJobs:   acceptJobs,
	}
	runnerId := config.GlobalConfig.Server.RunnerId
	runnerToken := config.GlobalConfig.Server.Token
	aliveResponse, err := serverRequest(runnerId, runnerToken, &aliveRequest, &alive.httpclient)
	if err != nil {
		return nil, err
	}
	return aliveResponse, nil
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

func serverRequest(runnerId string, runnerToken string, aliveRequest *AliveRequest, client *http.Client) (*AliveResponse, error) {
	address := config.GlobalConfig.Server.Address
	jsonBody, err := json.Marshal(*aliveRequest)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	request.Header.Add("RUNNER_ID", runnerId)
	request.Header.Add("RUNNER_TOKEN", runnerToken)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		body, _ := readBody(response)
		errorMsg := fmt.Sprintf("request alive server error,httpcode:%d ,body:%s", response.StatusCode, body)
		return nil, errors.New(errorMsg)
	}
	body, err := readBody(response)
	if err != nil {
		return nil, err
	}
	var aliveResponse AliveResponse
	err = json.Unmarshal(body, &aliveResponse)
	if err != nil {
		return nil, err
	}
	return &aliveResponse, nil
}

func readBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseBodyBytes, nil
}

type AliveResponse struct {
	NewJobs []work.NewJob `json:"newJobs"`
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
