package work

import (
	"client/config"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

type NewWork struct {
	PipeId     string           `json:"pipeId"` //流水线ID
	StepId     string           `json:"stepId"` //Step ID
	Type       int              `json:"type"`   //Work的类型,比如Command/Deploy类型
	WorkConfig *json.RawMessage `json:"workConfig"`
}

const (
	CommandType       = 1
	DeployType        = 2
	DockerCommandType = 3
)

type WorkerStarter struct {
	source  []Source
	worker  Worker
	workDir *WorkDir
	stepLog JobLog
}

func NewWorkerStarter(sources []Source, newJob *NewJob) (*WorkerStarter, error) {
	clientWorkDir := config.GlobalConfig.Local.ClientWorkDir

	stepWorkDirPath := clientWorkDir + "/" + newJob.PipeRunningId + "/" + newJob.JobRunningId
	_, err := os.Stat(stepWorkDirPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(stepWorkDirPath, fs.ModePerm)
	}
	var mainSourceName string
	if sources == nil || len(sources) == 0 {
		mainSourceName = ""
	} else {
		for _, source := range sources {
			if source.IsMainSource {
				mainSourceName = source.ProjectName
			}
		}
	}
	workDir, err := NewWorkDir(stepWorkDirPath, mainSourceName)
	if err != nil {
		return nil, err
	}
	stepLog := NewStepLog()
	worker, err := newWorker(newJob, workDir, &stepLog)
	err = workDir.CleanWorkDir()
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

func (starter *WorkerStarter) RunStarter() error {
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

func (starter *WorkerStarter) StepLog() JobLog {
	return starter.stepLog
}

type Worker interface {
	// Run the worker
	Run() error

	// Stop the worker
	Stop() error
}

func newWorker(newWork *NewWork, workDir *WorkDir, stepLog *JobLog) (Worker, error) {
	switch newWork.Type {
	case CommandType:
		body := newWork.WorkConfig
		return newCommandWorker(stepLog, body, workDir)
	case DeployType:
		// TODO Not support yet
		return nil, errors.New("DeployType not support yet")
	}
	return nil, errors.New("not support work type")
}

func handleResources(sources []Source, workDir *WorkDir, stepLog *JobLog) error {
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
