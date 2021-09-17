package work

import "log"

type CommandWorker struct {
	cmds       []string
	envs       []string
	pipeJobDir *PipeJobDir
	jobLog     *JobLog
	stopChan   chan bool
}

func newCommandWorker(stepLog *JobLog, cmdJob *CmdJobDto, pipeJobDir *PipeJobDir) (*CommandWorker, error) {
	cmdWorker := CommandWorker{
		cmds:       cmdJob.Cmds,
		envs:       cmdJob.Envs,
		pipeJobDir: pipeJobDir,
		jobLog:     stepLog,
	}
	return &cmdWorker, nil
}

func (cmdWorker *CommandWorker) Run() error {
	log.Println("Running a command type job")
	for _, cmd := range cmdWorker.cmds {
		actionLog := cmdWorker.jobLog.NewAction("Execute command: " + cmd)
		// TODO MainSourceDir is nil
		exec := NewExec(cmdWorker.pipeJobDir.MainSourceDir(), &actionLog, cmdWorker.envs, 10000, true, &cmdWorker.stopChan)
		err := exec.ExecShell(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmdWorker *CommandWorker) Stop() error {
	cmdWorker.stopChan <- true
	return nil
}
