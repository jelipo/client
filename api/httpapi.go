package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
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

func NewRunnerHttpApi(address string, runnerId string, runnerToken string) RunnerHttpApi {
	return RunnerHttpApi{
		runnerId:    runnerId,
		runnerToken: runnerToken,
		address:     address,
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
		errorMsg := fmt.Sprintf("request dsa server error,httpcode:%d ,body:%s", httpResponse.StatusCode, body)
		return errors.New(errorMsg)
	}
	body, err := readBody(httpResponse)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, responseInterface)
	if err != nil {
		logrus.Error("get response failed")
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
