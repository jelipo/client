package work

//
//import "os"
//
//type WorkDir struct {
//	StepWorkDir      string
//	ResourcesWorkDir string
//	MainWorkDir      string
//	TempWorkDir      string
//	MainSourceDir    string
//}
//
//func NewWorkDir(stepWorkDir string, mainSourceName string) (*WorkDir, error) {
//	resourcesWorkDir := stepWorkDir + "/resources"
//	mainWorkDir := stepWorkDir + "/work"
//	tempWorkDir := stepWorkDir + "/temp"
//	mainSourceDir := mainWorkDir + "/" + mainSourceName
//	err := mkdirDirs([]string{resourcesWorkDir, mainWorkDir, tempWorkDir})
//	if err != nil {
//		return nil, err
//	}
//	return &WorkDir{
//		StepWorkDir:      stepWorkDir,
//		ResourcesWorkDir: resourcesWorkDir,
//		MainWorkDir:      mainWorkDir,
//		TempWorkDir:      tempWorkDir,
//		MainSourceDir:    mainSourceDir,
//	}, nil
//}
//
//func (workDir WorkDir) ProjectMainWorkDir(projectName string) string {
//	return workDir.MainWorkDir + "/" + projectName
//}
//
//func (workDir WorkDir) CleanTempDir() error {
//	tempDir := workDir.TempWorkDir
//	err := os.RemoveAll(tempDir)
//	if err != nil {
//		return nil
//	}
//	err = os.MkdirAll(workDir.TempWorkDir, os.ModePerm)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (workDir WorkDir) CleanWorkDir() error {
//	mainWorkDir := workDir.MainWorkDir
//	err := os.RemoveAll(mainWorkDir)
//	if err != nil {
//		return nil
//	}
//	err = os.MkdirAll(workDir.MainWorkDir, os.ModePerm)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func mkdirDirs(dirs []string) error {
//	for _, dir := range dirs {
//		err := os.MkdirAll(dir, os.ModePerm)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
