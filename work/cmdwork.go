package work

import (
	"client/api"
	"github.com/sirupsen/logrus"
)

type CommandWorker struct {
	cmds       []string
	cmdEnvs    []string
	pipeEnvs   []api.PipeEnv
	pipeJobDir *PipeJobDir
	jobLog     *JobLog
	stopChan   chan bool
}

func newCommandWorker(stepLog *JobLog, cmdJob *api.CmdJobDto, pipeJobDir *PipeJobDir) (*CommandWorker, error) {
	cmdWorker := CommandWorker{
		cmds:       cmdJob.Cmds,
		cmdEnvs:    cmdJob.Envs,
		pipeEnvs:   nil,
		pipeJobDir: pipeJobDir,
		jobLog:     stepLog,
		stopChan:   make(chan bool, 100),
	}
	return &cmdWorker, nil
}

func (cmdWorker *CommandWorker) Run() error {
	logrus.Info("Running a command type job")
	for _, cmd := range cmdWorker.cmds {
		actionLog := cmdWorker.jobLog.NewAction("Execute command: " + cmd)
		// TODO MainSourceDir is nil
		envs := cmdWorker.cmdEnvs
		strings := changeEnvs(cmdWorker.pipeEnvs)
		envs = append(envs, strings...)
		exec := NewExec(cmdWorker.pipeJobDir.MainSourceDir(), &actionLog, cmdWorker.cmdEnvs, 5*60*1000, true, &cmdWorker.stopChan)
		err := exec.ExecShell(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func changeEnvs(pipeEnvs []api.PipeEnv) []string {
	var envsStr []string
	for _, env := range pipeEnvs {
		var envStr = env.Name + "=" + env.Value
		envsStr = append(envsStr, envStr)
	}
	return envsStr
}

func (cmdWorker *CommandWorker) Stop() error {
	cmdWorker.stopChan <- true
	return nil
}
