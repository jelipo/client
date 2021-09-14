package work

import (
	"errors"
)

type Source struct {
	SourceId        string          `json:"sourceId"`
	SourceType      string          `json:"sourceType"`
	UseCache        bool            `json:"useCache"`
	ProjectName     string          `json:"projectName"`
	GitSourceConfig GitSourceConfig `json:"gitSourceConfig"`
}

const (
	OutsideGit = "OUTSIDE_GIT"
	HttpFile   = "HTTP_FILE"
)

type Handler interface {
	// HandleSource download the source
	// return the resource path
	HandleSource() (*string, error)
}

func NewSourceHandler(source *Source, resourcesWorkDor string, stepLog *JobLog) (Handler, error) {
	switch source.SourceType {
	case OutsideGit:
		return NewGitSourceHandler(resourcesWorkDor, source.ProjectName, source.GitSourceConfig, stepLog)
	case HttpFile:
		return nil, errors.New("not support HttpDownloadType yet")
	}
	return nil, errors.New("not support yet")
}
