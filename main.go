package main

import (
	"client/work"
	"fmt"
	"time"
)

func main() {
	var workdir = "/home/cao/client_temp"
	stepLog := work.NewStepLog()
	var cmdStr = "rustup update"
	var env = make([]string, 0)
	actionLog := stepLog.NewAction("execute: " + cmdStr)
	exec := work.NewExec(workdir, &actionLog, env, 0)

	go printLog(&stepLog)

	err := exec.ExecShell("rustup update")
	if err != nil {
		fmt.Println("cmd error:" + err.Error())
	}
	fmt.Println("done")
	time.Sleep(time.Duration(99999999) * time.Millisecond)
}

func printLog(stepLog *work.StepLog) {
	for true {
		time.Sleep(time.Duration(500) * time.Millisecond)
		logs := stepLog.GetLogs(100)
		if logs == nil {
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}
		for _, log := range logs {
			fmt.Print(log.LogBody)
		}
	}

}
