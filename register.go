package main

import (
	"client/api"
	"client/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/user"
)

func RegisterToServer(initValue *InitValue) error {
	server := config.Server{}
	err := toServer(&server, initValue)
	if err != nil {
		return err
	}
	local := config.Local{}
	err = toLocal(&local, initValue)
	if err != nil {
		return err
	}
	newConfig := config.Config{
		Server: server,
		Local:  local,
	}
	err = config.WriteConfigFile(&newConfig, initValue.ConfigPath)
	if err != nil {
		return err
	}
	return nil
}

func toServer(server *config.Server, initValue *InitValue) error {
	var host string
	var err error
	if len(initValue.ServerHost) == 0 {
		host, err = scanln("please enter dps server host.(such as 'https://dps.daocloud.com')")
		if err != nil {
			return err
		}
	} else {
		host = initValue.ServerHost
	}
	var serverPath string
	if len(initValue.ServerPath) == 0 {
		serverPath = "/cicdengine/api/v1/runner"
		logrus.Info("use default dps server path: " + serverPath)
	} else {
		serverPath = initValue.ServerPath
	}
	address := host + serverPath
	tempHttpApi := api.NewRunnerHttpApi(address, "", "")
	var oneTimeToken = initValue.OneTimeToken
	if len(oneTimeToken) == 0 {
		oneTimeToken, err = scanln("please enter dps server one time token")
		if err != nil {
			return err
		}
	}
	response, err := tempHttpApi.RegisterToServer(oneTimeToken)
	if err != nil {
		return err
	}
	server.Token = *response.Token
	server.RunnerId = *response.RunnerId
	server.MaxJobNum = *response.MaxJobNum
	server.Address = address
	return nil
}

func scanln(tip string) (string, error) {
	fmt.Println(tip)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	fmt.Println("your input \"" + input + "\"")
	return input, nil
}

func toLocal(local *config.Local, initValue *InitValue) error {
	var workDir = initValue.WorkDir
	currUser, err := user.Current()
	if err != nil {
		return err
	}
	if len(workDir) == 0 {
		workDir = currUser.HomeDir + "/dps_runner/work"
	}
	local.ClientWorkDir = workDir
	return nil
}
