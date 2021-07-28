package work

import (
	"client/config"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

type NewWork struct {
	PipeSoleId string //执行流水线的唯一ID
	PipeId     string //流水线ID
	StepId     string //Step ID
	Type       int    //Work的类型,比如Command/Deploy类型
	WorkBody   *json.RawMessage
}

const (
	CommandType = 1
	DeployType  = 2
)

type WorkerStarter struct {
	source  []Source
	worker  Worker
	workDir *WorkDir
	stepLog StepLog
}

func NewWorkerStarter(sources []Source, newWork *NewWork) (*WorkerStarter, error) {
	clientWorkDir := config.GlobalConfig.Local.ClientWorkDir
	stepWorkDirPath := clientWorkDir + "/" + newWork.PipeId + "/" + newWork.StepId
	_, err := os.Stat(stepWorkDirPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(stepWorkDirPath, fs.ModeDir)
	}
	workDir, err := NewWorkDir(stepWorkDirPath)
	if err != nil {
		return nil, err
	}
	stepLog := NewStepLog()
	worker, err := newWorker(newWork, workDir, &stepLog)
	if err != nil {
		return nil, err
	}
	return &WorkerStarter{
		source:  sources,
		worker:  worker,
		workDir: workDir,
		stepLog: stepLog,
	}, nil
}

func (starter *WorkerStarter) Run() error {
	// Handle the resources
	err := handleResources(starter.source, starter.workDir, &starter.stepLog)
	if err != nil {
		return err
	}
	err = starter.worker.Run()
	if err != nil {
		return err
	}
	return nil
}

func (starter *WorkerStarter) StepLog() StepLog {
	return starter.stepLog
}

type Worker interface {
	// Run the worker
	Run() error
}

func newWorker(newWork *NewWork, workDir *WorkDir, stepLog *StepLog) (Worker, error) {
	switch newWork.Type {
	case CommandType:
		body := newWork.WorkBody
		return newCommandWorker(stepLog, body, workDir)
	case DeployType:
		// TODO Not support yet
		return nil, errors.New("DeployType not support yet")
	}
	return nil, errors.New("not support work type")
}

func handleResources(sources []Source, workDir *WorkDir, stepLog *StepLog) error {
	if sources != nil && len(sources) != 0 {
		for _, resource := range sources {
			//判断是否需要缓存resource
			var resourcesWorkDir string
			if resource.UseCache {
				resourcesWorkDir = workDir.ResourcesWorkDir
			} else {
				resourcesWorkDir = workDir.TempWorkDir
			}
			handler, err := NewSourceHandler(&resource, resourcesWorkDir, stepLog)
			if err != nil {
				return err
			}
			resourceProjectPath, err := handler.HandleSource()
			if err != nil {
				return err
			}
			err = os.Rename(*resourceProjectPath, workDir.ProjectMainWorkDir(resource.ProjectName))
			if err != nil {
				return err
			}
		}
	}
	err := workDir.CleanTempDir()
	if err != nil {
		return err
	}
	return nil
}

func cleanTemp() {

}
