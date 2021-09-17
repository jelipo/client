package work

import (
	"errors"
)

const (
	OutsideGit = "OUTSIDE_GIT"
	HttpFile   = "HTTP_FILE"
)

type Handler interface {
	// StartHandleSource HandleSource download the source
	// return the source path
	StartHandleSource() error
}

func NewSourceHandler(source *Source, resourcesWorkDir string, stepLog *JobLog) (Handler, error) {
	switch source.SourceType {
	case OutsideGit:
		return NewGitSourceHandler(resourcesWorkDir, source.ProjectName, &source.GitSourceConfig, stepLog)
	case HttpFile:
		return nil, errors.New("not support HttpDownloadType")
	}
	return nil, errors.New("not support")
}
