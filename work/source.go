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
	// Git仓库绝对文件路径
	repoAbPath      string
	gitSourceConfig *GitSourceConfig
}

func NewGitSourceHandler(repoAbPath string, gitConfig *GitSourceConfig) GitSourceHandler {
	return GitSourceHandler{repoAbPath: repoAbPath, gitSourceConfig: gitConfig}
}

func (gitHandler GitSourceHandler) HandleSource() error {
	// TODO Pull git repo

	return nil
}
