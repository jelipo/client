package main

import (
	"client/config"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	initValue, err := InitContext()
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
	logrus.Info("DPS Runner")
	// Find config file,
	exited, err := ConfigExist(initValue.ConfigPath)
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
	// if not exited try to register, and generate config file
	if !exited {
		logrus.Info("can not found config file,runner need to regist")
		err := RegisterToServer(initValue)
		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}
	}
	// Read config
	err = config.InitConfig(initValue.ConfigPath)
	if err != nil {
		return
	}
	// Listening
	pipeManager := NewRunnerClient()
	pipeManager.run()
}
