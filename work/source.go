package work

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

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
	Branch   string
	CommitId string
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

func (gitHandler GitSourceHandler) HandleSource() (*string, error) {
	repo, err := initRepo(gitHandler.gitRepoDir)
	if err != nil {
		return nil, err
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}
	gitAddr := gitHandler.gitSourceConfig.GitAddress
	if remote == nil {
		_, err := repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{gitAddr},
		})
		if err != nil {
			return nil, err
		}
	} else {
		remote.Config().URLs = []string{gitAddr}
	}
	err = repo.Fetch(&git.FetchOptions{RemoteName: "origin", Depth: 1})
	if err != nil {
		return nil, err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	err = worktree.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(gitHandler.gitSourceConfig.CommitId)})
	if err != nil {
		return nil, err
	}
	return &gitHandler.gitRepoDir, nil
}

func initRepo(gitRepoDir string) (*git.Repository, error) {
	var repo *git.Repository
	_, err := os.Stat(gitRepoDir)
	if err != nil && os.IsExist(err) {
		repo, err = git.PlainOpen(gitRepoDir)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		repo, err = git.PlainInit(gitRepoDir, false)
		if err != nil {
			return nil, err
		}
	}
	return repo, nil
}
