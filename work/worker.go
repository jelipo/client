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
	sources    []Source
	WorkBody   *json.RawMessage
}

const (
	CommandType = 1
	DeployType  = 2
)

type Worker interface {
	Start() error
}

func NewWorker(newWork *NewWork) (Worker, error) {
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
	switch newWork.Type {
	case CommandType:
		body := newWork.WorkBody
		return newCommandWorker(&stepLog, body, workDir)
	case DeployType:
		// TODO Not support yet
		return nil, errors.New("not support yet")
	}
	// Handle the resources
	err = handleResources(newWork.sources, workDir, &stepLog)
	if err != nil {
		return nil, err
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
