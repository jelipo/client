package work

import (
	"client/api"
	"errors"
)

type Handler interface {
	// StartHandleSource HandleSource download the source
	// return the source path
	StartHandleSource() error
}

func NewSourceHandler(source *api.Source, resourcesWorkDir string, stepLog *JobLog) (Handler, error) {
	switch source.SourceType {
	case api.OutsideGit:
		return NewGitSourceHandler(resourcesWorkDir, source.ProjectName, &source.GitSourceConfig, stepLog)
	case api.HttpFile:
		return nil, errors.New("not support HttpDownloadType")
	}
	return nil, errors.New("not support")
}
