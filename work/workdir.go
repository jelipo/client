package work

import "os"

type WorkDir struct {
	StepWorkDir      string
	ResourcesWorkDir string
	MainWorkDir      string
	TempWorkDir      string
}

func NewWorkDir(stepWorkDir string) (*WorkDir, error) {
	resourcesWorkDir := stepWorkDir + "/resources"
	mainWorkDir := stepWorkDir + "/work"
	tempWorkDir := stepWorkDir + "/temp"
	err := mkdirDirs([]string{resourcesWorkDir, mainWorkDir, tempWorkDir})
	if err != nil {
		return nil, err
	}
	return &WorkDir{
		StepWorkDir:      stepWorkDir,
		ResourcesWorkDir: resourcesWorkDir,
		MainWorkDir:      mainWorkDir,
		TempWorkDir:      tempWorkDir,
	}, nil
}

func (workDir WorkDir) ProjectMainWorkDir(projectName string) string {
	return workDir.MainWorkDir + "/" + projectName
}

func (workDir WorkDir) CleanTempDir() error {
	tempDir := workDir.TempWorkDir
	err := os.Remove(tempDir)
	if err != nil {
		return nil
	}
	err = os.MkdirAll(workDir.TempWorkDir, os.ModeDir)
	if err != nil {
		return err
	}
	return nil
}

func mkdirDirs(dirs []string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}
