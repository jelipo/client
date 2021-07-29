package main

import "client/work"

type RunnerAlive struct {
	serverAddress string
	token         string
}

func NewRunnerAlive(serverAddr string, token string) RunnerAlive {
	return RunnerAlive{
		serverAddress: serverAddr,
		token:         token,
	}
}

func (alive *RunnerAlive) alive() AliveResponse {
	// TODO http
	return AliveResponse{}
}

type NewJob struct {
	Sources []work.Source `json:"sources"`
	NewWork work.NewWork  `json:"newWork"`
}

type AliveResponse struct {
	NewJobs []NewJob `json:"newJobs"`
}
