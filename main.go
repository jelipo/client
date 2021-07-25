package main

import (
	"client/work"
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"strings"
	"time"
)

func main() {
	fs := osfs.New("/home/cao/client_temp/git")\
	fs.

	var workdir = "/home/cao/client_temp"
	stepLog := work.NewStepLog()
	var cmdStr = "export TEST=1"
	var env = make([]string, 0)
	actionLog := stepLog.NewAction("execute: " + cmdStr + "\n")
	exec := work.NewExec(workdir, &actionLog, env, 0)

	go printLog(&stepLog)

	err := exec.ExecShell(cmdStr)
	if err != nil {
		fmt.Println("cmd error:" + err.Error())
	}
	fmt.Println("done")
	envs := exec.GetEnvs()
	for _, env := range *envs {
		fmt.Println(env)
	}
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
			r1 := strings.ReplaceAll(log.LogBody, "\r", "\n")
			fmt.Print(r1)
			//fmt.Print(log.LogBody)
		}
	}
}
