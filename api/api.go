package api

import "runtime"

type RegisterRequest struct {
	OneTimeToken string `json:"oneTimeToken"`
	Os           string `json:"os"`
	Arch         string `json:"arch"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

func (api *RunnerHttpApi) register(oneTimeToken string) (*RegisterResponse, error) {
	registerRequest := RegisterRequest{
		OneTimeToken: oneTimeToken,
		Os:           runtime.GOOS,
		Arch:         runtime.GOARCH,
	}
	var response RegisterResponse
	err := api.doHttp("POST", api.address+"/register", registerRequest, response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
