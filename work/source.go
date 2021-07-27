package work

import (
	"encoding/json"
	"errors"
)

type Source struct {
	Type         int
	UseCache     bool //是否使用缓存
	ProjectName  string
	SourceConfig *json.RawMessage
}

const (
	GitSourceType    = 1
	HttpDownloadType = 2
)

type Handler interface {
	// HandleSource download the source
	// return the resource path
	HandleSource() (*string, error)
}

func NewSourceHandler(source *Source, stepWorkDir *WorkDir, stepLog *StepLog) (Handler, error) {
	resourcesDir := stepWorkDir.ResourcesWorkDir
	switch source.Type {
	case GitSourceType:
		return NewGitSourceHandler(resourcesDir, source.ProjectName, source.SourceConfig, stepLog)
	case HttpDownloadType:
		return nil, errors.New("not support HttpDownloadType yet")
	}
	return nil, errors.New("not support yet")
}
