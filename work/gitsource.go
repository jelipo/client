package work

import (
	"client/api"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"strings"
)

type GitSourceHandler struct {
	gitSourceConfig *api.GitSourceConfig
	repoName        string
	gitRepoDir      string
	jobLog          *JobLog
	sourceName      string
}

func NewGitSourceHandler(sourceDir string, sourceName string, gitSourceConfig *api.GitSourceConfig, jobLog *JobLog) (*GitSourceHandler, error) {
	return &GitSourceHandler{
		gitSourceConfig: gitSourceConfig,
		gitRepoDir:      sourceDir,
		jobLog:          jobLog,
		sourceName:      sourceName,
	}, nil
}

func (gitHandler *GitSourceHandler) StartHandleSource() error {
	actionLog := gitHandler.jobLog.NewAction("Get the git source file. sourceName:" + gitHandler.sourceName)
	_ = os.MkdirAll(gitHandler.gitRepoDir, os.ModePerm)
	if len(gitHandler.gitSourceConfig.CommitId) != 0 {
		repo, err := gitInitRepo(gitHandler.gitRepoDir)
		if err != nil {
			actionLog.AddExecLog("Git source checkout failed.")
			return err
		}
		err = fetchGitFile(repo, gitHandler.gitSourceConfig, &actionLog)
		if err != nil {
			return err
		}
	} else {
		err := gitPlainClone(gitHandler.gitRepoDir, gitHandler.gitSourceConfig, &actionLog)
		if err != nil {
			actionLog.AddExecLog("Git source plain clone failed.")
			return err
		}
	}
	actionLog.AddExecLog("Git source download success.")
	return nil
}

func gitInitRepo(gitRepoDir string) (*git.Repository, error) {
	var repo *git.Repository
	_, err := os.Stat(gitRepoDir)
	repo, err = git.PlainInit(gitRepoDir, false)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func fetchGitFile(repo *git.Repository, gitSourceConfig *api.GitSourceConfig, actionLog *ActionLog) error {
	remotes, _ := repo.Remotes()
	remote, _ := repo.Remote("origin")
	if len(remotes) == 0 || remote == nil {
		_, err := repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{gitSourceConfig.GitAddress},
		})
		if err != nil {
			return err
		}
	} else {
		remote.Config().URLs = []string{gitSourceConfig.GitAddress}
	}
	//Set Git Auth
	auth, err := gitAuth(gitSourceConfig)
	if err != nil {
		return err
	}
	fetchErr := repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   actionLog,
		Auth:       auth,
	})

	if fetchErr != nil && !strings.Contains(fetchErr.Error(), "already up-to-date") {
		return fetchErr
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	//
	var checkoutOptions git.CheckoutOptions
	checkoutOptions = git.CheckoutOptions{
		Hash:  plumbing.NewHash(gitSourceConfig.CommitId),
		Force: true,
	}
	err = worktree.Checkout(&checkoutOptions)
	if err != nil {
		return err
	}
	return nil
}

func gitPlainClone(path string, gitSourceConfig *api.GitSourceConfig, actionLog *ActionLog) error {
	auth, err := gitAuth(gitSourceConfig)
	if err != nil {
		return err
	}
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:           gitSourceConfig.GitAddress,
		Auth:          auth,
		ReferenceName: plumbing.NewBranchReferenceName(gitSourceConfig.Branch),
		SingleBranch:  true,
		NoCheckout:    false,
		Depth:         1,
		Progress:      actionLog,
	})
	if err != nil {
		return err
	}
	return nil
}

func gitAuth(config *api.GitSourceConfig) (transport.AuthMethod, error) {
	switch config.AuthType {
	case api.GitAuthPassword:
		auth := http.BasicAuth{Username: config.AuthUsername, Password: config.AuthPassword}
		return &auth, nil
	case api.GitAuthPublicKeyFile:
		publicKey, err := ssh.NewPublicKeysFromFile("git", config.AuthPublicKeyPath, config.AuthPassword)
		if err != nil {
			return nil, err
		}
		return publicKey, err
	case api.GitAuthPublicKeyStr:
		publicKey, err := ssh.NewPublicKeys("git", []byte(config.AuthPublicKeyStr), config.AuthPassword)
		if err != nil {
			return nil, err
		}
		return publicKey, err
	case api.NoAuth:
		auth := http.BasicAuth{Username: config.AuthUsername, Password: config.AuthPassword}
		return &auth, nil
	}
	return nil, errors.New("unknown auth type")
}
