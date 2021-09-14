package work

import (
	"client/work"
	"os"
)

type PipeJobDir struct {
	RunningPipeDir string
	WorkDir        string
	TempWorkDir    string
	MainSourceId   string
	SourceIdDirMap map[string]string
}

func NewPipeJobDir(runnerClientWorkDir string, pipeRunningId string, jobRunningId string, mainSourceId string, sources []work.Source) (*PipeJobDir, error) {
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
		WorkDir:        workDir,
		TempWorkDir:    tempWorkDir,
		MainSourceId:   mainSourceId,
		SourceIdDirMap: sourceIdDirMap,
	}, nil
}

func (workDir PipeJobDir) SourceDir(sourceId string) string {
	return workDir.SourceDir(sourceId)
}

func (workDir PipeJobDir) CleanTempDir() error {
	return cleanDir(workDir.TempWorkDir)
}

func (workDir PipeJobDir) CleanWorkDir() error {
	return cleanDir(workDir.WorkDir)
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
