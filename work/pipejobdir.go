package work

import (
	"client/api"
	"os"
)

type PipeJobDir struct {
	RunningPipeDir string
	RunningJobDir  string
	WorkDir        string
	TempWorkDir    string
	MainSourceId   string
	SourceIdDirMap map[string]string
}

func NewPipeJobDir(runnerClientWorkDir string, pipeRunningId string, jobRunningId string, mainSourceId string, sources []api.Source) (*PipeJobDir, error) {
	runningPipeDir := runnerClientWorkDir + "/" + pipeRunningId
	runningJobDir := runningPipeDir + "/" + jobRunningId
	workDir := runningJobDir + "/work"
	tempWorkDir := runningJobDir + "/temp"
	mkDirs := []string{workDir, tempWorkDir}
	sourceIdDirMap := make(map[string]string)
	if len(mainSourceId) != 0 && len(sources) != 0 {
		for _, source := range sources {
			sourceDir := workDir + "/" + source.ProjectName
			mkDirs = append(mkDirs, sourceDir)
			sourceIdDirMap[source.SourceId] = sourceDir
		}
	}
	err := mkdirDirs(mkDirs)
	if err != nil {
		return nil, err
	}
	return &PipeJobDir{
		RunningPipeDir: runningPipeDir,
		RunningJobDir:  runningJobDir,
		WorkDir:        workDir,
		TempWorkDir:    tempWorkDir,
		MainSourceId:   mainSourceId,
		SourceIdDirMap: sourceIdDirMap,
	}, nil
}

func (workDir PipeJobDir) SourceDir(sourceId string) string {
	return workDir.SourceIdDirMap[sourceId]
}

func (workDir PipeJobDir) MainSourceDir() string {
	return workDir.SourceDir(workDir.MainSourceId)
}

func (workDir PipeJobDir) CleanTempDir() error {
	return cleanDir(workDir.TempWorkDir)
}

func (workDir PipeJobDir) CleanJobWorkDir() error {
	return cleanDir(workDir.WorkDir)
}

func (workDir PipeJobDir) CleanRunningJobDir() error {
	return cleanDir(workDir.RunningJobDir)
}

func cleanDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return nil
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func mkdirDirs(dirs []string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
