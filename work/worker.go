package work

import (
	"client/api"
	"client/config"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
)

type NewWork struct {
	PipeId     string           `json:"pipeId"` //流水线ID
	StepId     string           `json:"stepId"` //Step ID
	Type       int              `json:"type"`   //Work的类型,比如Command/Deploy类型
	WorkConfig *json.RawMessage `json:"workConfig"`
}

type JobStarter struct {
	sources      []api.Source
	worker       JobWorker
	pipeJobDir   *PipeJobDir
	jobLog       *JobLog
	jobRunningId string
	mainSourceId string
}

func NewJobStarter(sources []api.Source, newJob *api.NewJob) (*JobStarter, error) {
	clientWorkDir := config.GlobalConfig.Local.ClientWorkDir

	pipeJobDir, err := NewPipeJobDir(clientWorkDir, newJob.PipeRunningId, newJob.JobRunningId, newJob.MainSourceId, newJob.Sources)
	if err != nil {
		return nil, err
	}
	stepLog := NewJobLog()
	worker, err := newWorker(newJob, pipeJobDir, &stepLog)
	err = pipeJobDir.CleanJobWorkDir()
	if err != nil {
		return nil, err
	}
	return &JobStarter{
		sources:      sources,
		worker:       worker,
		pipeJobDir:   pipeJobDir,
		jobLog:       &stepLog,
		jobRunningId: newJob.JobRunningId,
		mainSourceId: newJob.MainSourceId,
	}, nil
}

func (starter *JobStarter) RunStarter() error {
	// Handle the resources
	sourceResult, err := handleResources(starter.sources, starter.pipeJobDir, starter.jobLog, starter.mainSourceId)
	if err != nil {
		return err
	}
	logrus.Info("Running job,jobRunningId:" + starter.jobRunningId)
	return starter.worker.Run(sourceResult)
}

func (starter *JobStarter) JobLog() *JobLog {
	return starter.jobLog
}

type JobWorker interface {
	// Run the worker
	Run(result *SourceResult) error

	// Stop the worker
	Stop() error
}

func newWorker(newJob *api.NewJob, pipeJobDir *PipeJobDir, jobLog *JobLog) (JobWorker, error) {
	switch newJob.JobType {
	case api.CommandType:
		return newCommandWorker(jobLog, &newJob.CmdJobDto, pipeJobDir)
	case api.DeployType:
		// TODO Not support yet
		return nil, errors.New("DeployType not support yet")
	}
	return nil, errors.New("not support work type")
}

func handleResources(sources []api.Source, pipeJobDir *PipeJobDir, jobLog *JobLog, mainSourceId string) (*SourceResult, error) {
	if len(sources) != 0 {
		for _, source := range sources {
			//判断是否需要缓存resource
			var resourcesWorkDir = pipeJobDir.SourceDir(source.SourceId)
			isMainSource := source.SourceId == mainSourceId
			handler, err := NewSourceHandler(&source, resourcesWorkDir, jobLog, isMainSource)
			if err != nil {
				return nil, err
			}
			return handler.StartHandleSource()
		}
	}
	return nil, nil
}
