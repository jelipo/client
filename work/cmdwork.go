package work

import (
	"encoding/json"
	"fmt"
)

type CommandWork struct {
	Sources []Source `json:"sources"`
	Cmds    []string `json:"cmds"`
	Envs    []string `json:"envs"`
}

type CommandWorker struct {
	sources []Source
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
		sources: cmdWork.Sources,
		cmds:    cmdWork.Cmds,
		envs:    cmdWork.Envs,
		workDir: stepWorkDir,
		stepLog: stepLog,
	}
	return &cmdWorker, nil
}

// Start the worker
func (cmdWorker *CommandWorker) Start() error {
	err := handleResources(cmdWorker.sources, cmdWorker.workDir, cmdWorker.stepLog)
	if err != nil {
		return err
	}
	for _, cmd := range cmdWorker.cmds {
		actionLog := cmdWorker.stepLog.NewAction("Execute command: " + cmd)
		NewExec(cmdWorker.workDir.MainWorkDir, &actionLog, cmdWorker.envs, 10000, true)
	}
	return nil
}

func handleResources(sources []Source, stepWorkDir *WorkDir, stepLog *StepLog) error {
	if sources != nil && len(sources) != 0 {
		for _, resource := range sources {
			handler, err := NewSourceHandler(&resource, stepWorkDir, stepLog)
			if err != nil {
				return err
			}
			resourcesPath, err := handler.HandleSource()
			if err != nil {
				return err
			}
			// TODO move resource files to main work dir
			fmt.Println(resourcesPath)
		}
	}
	return nil
}
