package work

import (
	"encoding/json"
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

type GitSourceConfig struct {
	GitAddress        string `json:"gitAddress"`
	AuthType          int    `json:"authType"` //拉取git的身份验证方式,PublicKey or Password
	CommitId          string `json:"commitId"` //CommitHashId
	AuthUsername      string `json:"authUsername"`
	AuthPassword      string `json:"authPassword"`
	AuthPublicKeyStr  string `json:"authPublicKeyStr"`
	AuthPublicKeyPath string `json:"authPublicKeyPath"`
}

const (
	GitAuthPassword      = 1
	GitAuthPublicKeyStr  = 2
	GitAuthPublicKeyFile = 3
)

type GitSourceHandler struct {
	resourceDir     string
	gitSourceConfig *GitSourceConfig
	repoName        string
	gitRepoDir      string
	stepLog         *StepLog
}

func NewGitSourceHandler(resourceDir string, repoName string, sourceConfig *json.RawMessage, stepLog *StepLog) (*GitSourceHandler, error) {
	var gitResourceConfig GitSourceConfig
	err := json.Unmarshal(*sourceConfig, &gitResourceConfig)
	if err != nil {
		return nil, err
	}
	return &GitSourceHandler{
		resourceDir:     resourceDir,
		gitSourceConfig: &gitResourceConfig,
		repoName:        repoName,
		gitRepoDir:      resourceDir + "/" + repoName,
		stepLog:         stepLog,
	}, nil
}

func (gitHandler GitSourceHandler) HandleSource() (*string, error) {
	actionLog := gitHandler.stepLog.NewAction("Get the git resource")
	_ = os.MkdirAll(gitHandler.resourceDir, os.ModePerm)
	repo, err := gitInitRepo(gitHandler.gitRepoDir)
	if err != nil {
		return nil, err
	}
	gitAddr := gitHandler.gitSourceConfig.GitAddress
	remotes, _ := repo.Remotes()
	remote, _ := repo.Remote("origin")
	if len(remotes) == 0 || remote == nil {
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

	auth, err := gitAuth(gitHandler.gitSourceConfig)
	if err != nil {
		return nil, err
	}
	fetchErr := repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   actionLog,
		Auth:       auth,
	})
	if fetchErr != nil && !strings.Contains(fetchErr.Error(), "already up-to-date") {
		return nil, fetchErr
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(gitHandler.gitSourceConfig.CommitId),
		Force: true,
	})
	if err != nil {
		return nil, err
	}
	return &gitHandler.gitRepoDir, nil
}

func gitInitRepo(gitRepoDir string) (*git.Repository, error) {
	var repo *git.Repository
	_, err := os.Stat(gitRepoDir)
	if err != nil {
		if os.IsNotExist(err) {
			repo, err = git.PlainInit(gitRepoDir, false)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		repo, err = git.PlainOpen(gitRepoDir)
	}
	return repo, nil
}

func gitAuth(config *GitSourceConfig) (transport.AuthMethod, error) {
	switch config.AuthType {
	case GitAuthPassword:
		auth := http.BasicAuth{Username: config.AuthUsername, Password: config.AuthPassword}
		return &auth, nil
	case GitAuthPublicKeyFile:
		publicKey, err := ssh.NewPublicKeysFromFile("git", config.AuthPublicKeyPath, config.AuthPassword)
		if err != nil {
			return nil, err
		}
		return publicKey, err
	case GitAuthPublicKeyStr:
		publicKey, err := ssh.NewPublicKeys("git", []byte(config.AuthPublicKeyStr), config.AuthPassword)
		if err != nil {
			return nil, err
		}
		return publicKey, err
	}
	return nil, errors.New("unknown auth type")
}
