package work

type WorkDir struct {
	StepWorkDir      string
	ResourcesWorkDir string
	MainWorkDir      string
	TempWorkDir      string
}

func NewWorkDir(stepWorkDir string) WorkDir {
	return WorkDir{
		StepWorkDir:      stepWorkDir,
		ResourcesWorkDir: stepWorkDir + "/resources",
		MainWorkDir:      stepWorkDir + "/work",
		TempWorkDir:      stepWorkDir + "/temp",
	}
}
