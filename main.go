package main

import (
	"client/config"
	"client/work"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
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
		err := RegisterServer(initValue)
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

func register() error {
	return nil
}

func printLog(stepLog *work.JobLog) {
	for true {
		time.Sleep(time.Duration(500) * time.Millisecond)
		logs := stepLog.GetLogs(100)
		if logs == nil {

			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}
		for _, log := range logs {
			r1 := strings.ReplaceAll(log.LogBody, "\r", "\n")
			fmt.Print(r1)
			//fmt.Print(log.LogBody)
		}
	}
}
