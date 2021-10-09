package api

import (
	"client/work"
	"runtime"
)

type RegisterRequest struct {
	OneTimeToken *string `json:"oneTimeToken"`
	Os           *string `json:"os"`
	Arch         *string `json:"arch"`
}

type RegisterResponse struct {
	Token     *string `json:"token"`
	RunnerId  *string `json:"runnerId"`
	MaxJobNum *int    `json:"maxJobNum"`
}

func (api *RunnerHttpApi) RegisterToServer(oneTimeToken string) (*RegisterResponse, error) {
	os := runtime.GOOS
	arch := runtime.GOARCH
	registerRequest := RegisterRequest{
		OneTimeToken: &oneTimeToken,
		Os:           &os,
		Arch:         &arch,
	}
	var response RegisterResponse
	err := api.doHttp("POST", api.address+"/register", registerRequest, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type AliveResponse struct {
	NewJobs []work.NewJob `json:"newJobs"`
}

type AliveRequest struct {
	HostStatus   HostStatus         `json:"hostStatus"`
	RunnerStatus work.ManagerStatus `json:"runnerStatus"`
	JobsStatus   []JobsStatus       `json:"jobsStatus"`
	AcceptJobs   []string           `json:"acceptJobs"`
	DenyJobs     []string           `json:"denyJobs"`
}

type HostStatus struct {
	// TODO CPU/Memory/Disk info
}

type JobsStatus struct {
	JobRunningId   string         `json:"jobRunningId"`
	AtomLogs       []work.AtomLog `json:"atomLogs"`
	Finished       bool           `json:"finished"`
	FinishedStatus string         `json:"finishedStatus"`
}

func (api *RunnerHttpApi) AliveToServer(aliveRequest *AliveRequest) (*AliveResponse, error) {
	var response AliveResponse
	err := api.doHttp("POST", api.address+"/live", aliveRequest, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
