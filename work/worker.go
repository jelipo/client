package work

import (
	"client/config"
	"encoding/json"
	"errors"
	"log"
)

type NewWork struct {
	PipeId     string           `json:"pipeId"` //流水线ID
	StepId     string           `json:"stepId"` //Step ID
	Type       int              `json:"type"`   //Work的类型,比如Command/Deploy类型
	WorkConfig *json.RawMessage `json:"workConfig"`
}

const (
	CommandType       = "COMMAND"
	DeployType        = "DEPLOY"
	DockerCommandType = "DOCKER_COMMAND"
)

type JobStarter struct {
	source       []Source
	worker       JobWorker
	pipeJobDir   *PipeJobDir
	jobLog       *JobLog
	jobRunningId string
}

func NewJobStarter(sources []Source, newJob *NewJob) (*JobStarter, error) {
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
		source:       sources,
		worker:       worker,
		pipeJobDir:   pipeJobDir,
		jobLog:       &stepLog,
		jobRunningId: newJob.JobRunningId,
	}, nil
}

func (starter *JobStarter) RunStarter() error {
	// Handle the resources
	err := handleResources(starter.source, starter.pipeJobDir, starter.jobLog)
	if err != nil {
		return err
	}
	log.Println("Running job")
	err = starter.worker.Run()
	if err != nil {
		return err
	}
	return nil
}

func (starter *JobStarter) JobLog() *JobLog {
	return starter.jobLog
}

type JobWorker interface {
	// Run the worker
	Run() error

	// Stop the worker
	Stop() error
}

func newWorker(newJob *NewJob, pipeJobDir *PipeJobDir, jobLog *JobLog) (JobWorker, error) {
	switch newJob.JobType {
	case CommandType:
		return newCommandWorker(jobLog, &newJob.CmdJobDto, pipeJobDir)
	case DeployType:
		// TODO Not support yet
		return nil, errors.New("DeployType not support yet")
	}
	return nil, errors.New("not support work type")
}

func handleResources(sources []Source, pipeJobDir *PipeJobDir, jobLog *JobLog) error {
	if len(sources) != 0 {
		for _, source := range sources {
			//判断是否需要缓存resource
			var resourcesWorkDir = pipeJobDir.SourceDir(source.SourceId)
			handler, err := NewSourceHandler(&source, resourcesWorkDir, jobLog)
			if err != nil {
				return err
			}
			err = handler.StartHandleSource()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
