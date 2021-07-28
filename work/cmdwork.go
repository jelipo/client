package work

import (
	"encoding/json"
)

type CommandWork struct {
	Cmds []string `json:"cmds"`
	Envs []string `json:"envs"`
}

type CommandWorker struct {
	cmds    []string
	envs    []string
	workDir *WorkDir
	stepLog *StepLog
}

func newCommandWorker(stepLog *StepLog, workBody *json.RawMessage, stepWorkDir *WorkDir) (*CommandWorker, error) {
	var cmdWork CommandWork
	err := json.Unmarshal(*workBody, &cmdWork)
	if err != nil {
		return nil, err
	}
	cmdWorker := CommandWorker{
		cmds:    cmdWork.Cmds,
		envs:    cmdWork.Envs,
		workDir: stepWorkDir,
		stepLog: stepLog,
	}
	return &cmdWorker, nil
}

func (cmdWorker *CommandWorker) Run() error {
	for _, cmd := range cmdWorker.cmds {
		actionLog := cmdWorker.stepLog.NewAction("Execute command: " + cmd)
		NewExec(cmdWorker.workDir.MainWorkDir, &actionLog, cmdWorker.envs, 10000, true)
	}
	return nil
}
