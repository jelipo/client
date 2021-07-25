package work

type Source struct {
	Type         int
	SourceConfig interface{}
}

type GitSourceConfig struct {
	GitAddress string
	// 拉取git的身份验证方式
	AuthType int
	// 是否使用缓存
	UseCache bool
	// 拉取方式
	PullType int
}

type SourceHandler interface {
	// HandleSource download the source
	HandleSource() error
}

type GitSourceHandler struct {
	resourceDir     string
	gitSourceConfig *GitSourceConfig
	projectName     string
	gitRepoDir      string
	stepLog         *StepLog
}

func NewGitSourceHandler(resourceDir string, projectName string, gitConfig *GitSourceConfig, stepLog *StepLog) GitSourceHandler {
	return GitSourceHandler{
		resourceDir:     resourceDir,
		gitSourceConfig: gitConfig,
		projectName:     projectName,
		gitRepoDir:      resourceDir + "/" + projectName,
		stepLog:         stepLog,
	}
}

func (gitHandler GitSourceHandler) HandleSource() error {
	// TODO Pull git repo
	_ = gitHandler.stepLog.NewAction("Get the git resource : " + gitHandler.projectName)

	//exec := NewExec(gitHandler.gitRepoDir, &actionLog, make([]string, 0), 100000)
	//fs := osfs.New(gitHandler.gitRepoDir)
	//stat, err := fs.Stat(git.GitDirName)
	//if _, err := fs.Stat(git.GitDirName); err == nil {
	//	fs, err = fs.Chroot(git.GitDirName)
	//	CheckIfError(err)
	//}
	//
	//s := filesystem.NewStorageWithOptions(fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
	//git.Open(filesystem.NewStorage())
	return nil
}

func initRepo(actionLog *ActionLog, gitRepoDir string, executor Executor) error {
	var gitInitCmd = "git init " + gitRepoDir
	actionLog.AddSysLog(gitInitCmd)
	gitInitErr := executor.ExecShell(gitInitCmd)
	if gitInitErr != nil {
		return gitInitErr
	}
	return nil
}
