package main

import (
	"client/config"
	"client/work"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func main() {
	Init()
	logrus.Info("Hello World")
	err := config.IninConfig("/home/cao/go/client/config.json")
	if err != nil {
		return
	}
	pipeManager := NewRunnerClient()
	pipeManager.run()
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
