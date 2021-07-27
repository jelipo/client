package work

import (
	"client/config"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

type NewWork struct {
	PipeId   string
	StepId   string
	Type     int
	WorkBody *json.RawMessage
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
	workDir := NewWorkDir(stepWorkDirPath)
	stepLog := NewStepLog()
	switch newWork.Type {
	case CommandType:
		body := newWork.WorkBody
		return newCommandWorker(&stepLog, body, &workDir)
	case DeployType:
		// TODO Not support yet
	}
	return nil, errors.New("not support work type")
}
