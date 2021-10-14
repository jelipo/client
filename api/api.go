package api

import (
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
