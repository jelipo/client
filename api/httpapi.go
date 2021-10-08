package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RunnerHttpApi struct {
	runnerId    string
	runnerToken string
	address     string
	httpclient  http.Client
}

func NewRunnerHttpApi(address string) RunnerHttpApi {
	return RunnerHttpApi{
		address: address,
		httpclient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (api *RunnerHttpApi) doHttp(httpMethod string, url string, requestBody interface{}, responseInterface interface{}) error {
	jsonBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	httpRequest, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBodyBytes))
	if err != nil {
		return err
	}
	httpRequest.Header.Add("RUNNER_ID", api.runnerId)
	httpRequest.Header.Add("RUNNER_TOKEN", api.runnerToken)
	httpRequest.Header.Add("Content-Type", "application/json")
	httpResponse, err := api.httpclient.Do(httpRequest)
	if err != nil {
		return err
	}
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		body, _ := readBody(httpResponse)
		errorMsg := fmt.Sprintf("request alive server error,httpcode:%d ,body:%s", httpResponse.StatusCode, body)
		return errors.New(errorMsg)
	}
	body, err := readBody(httpResponse)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, responseInterface)
	if err != nil {
		return err
	}
	return nil
}

func readBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return responseBodyBytes, nil
}
