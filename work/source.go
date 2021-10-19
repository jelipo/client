package work

import (
	"client/api"
	"errors"
)

type Handler interface {
	// StartHandleSource HandleSource download the sources
	// return the sources path
	StartHandleSource() (*SourceResult, error)
}

func NewSourceHandler(source *api.Source, resourcesWorkDir string, stepLog *JobLog, isMainSource bool) (Handler, error) {
	switch source.SourceType {
	case api.OutsideGit:
		return NewGitSourceHandler(resourcesWorkDir, source.ProjectName, &source.GitSourceConfig, stepLog, isMainSource)
	case api.HttpFile:
		return nil, errors.New("not support HttpDownloadType")
	}
	return nil, errors.New("not support")
}

type SourceResult struct {
	SourceEnvs []SourceEnv
}

type SourceEnv struct {
	name  string
	value string
}
