package api

type AliveResponse struct {
	NewJobs []NewJob `json:"newJobs"`
}

type NewJob struct {
	//NewWork       work.NewWork  `json:"newWork"`
	JobRunningId  string    `json:"jobRunningId"`
	Sources       []Source  `json:"sources"`
	PipeId        string    `json:"pipeId"`
	StageId       string    `json:"stageId"`
	JobId         string    `json:"jobId"`
	PipeRunningId string    `json:"pipeRunningId"`
	MainSourceId  string    `json:"mainSourceId"`
	JobType       JobType   `json:"jobType"`
	PipeEnvs      []PipeEnv `json:"pipeEnvs"`
	CmdJobDto     CmdJobDto `json:"cmdJob"` // Exited when jobType is "COMMAND"
}

type JobType string

const (
	CommandType       = JobType("COMMAND")
	DeployType        = JobType("DEPLOY")
	DockerCommandType = JobType("DOCKER_COMMAND")
)

type PipeEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Source struct {
	SourceId        string          `json:"sourceId"`
	SourceType      SourceType      `json:"sourceType"`
	UseCache        bool            `json:"useCache"`
	ProjectName     string          `json:"projectName"`
	GitSourceConfig GitSourceConfig `json:"gitSourceConfig"`
}

type SourceType string

const (
	OutsideGit = SourceType("OUTSIDE_GIT")
	HttpFile   = SourceType("HTTP_FILE")
)

type GitSourceConfig struct {
	GitAddress        string      `json:"gitAddress"`
	AuthType          GitAuthType `json:"authType"` //拉取git的身份验证方式,PublicKey or Password
	Branch            string      `json:"branch"`
	CommitId          string      `json:"commitId"` //CommitHashId
	AuthUsername      string      `json:"authUsername"`
	AuthPassword      string      `json:"authPassword"`
	AuthPublicKeyStr  string      `json:"authPublicKeyStr"`
	AuthPublicKeyPath string      `json:"authPublicKeyPath"`
}

type GitAuthType string

const (
	NoAuth               = GitAuthType("NO_AUTH")
	GitAuthPassword      = GitAuthType("PASSWORD")
	GitAuthPublicKeyStr  = GitAuthType("PUBLIC_KEY_STR")
	GitAuthPublicKeyFile = GitAuthType("PUBLIC_KEY_FILE")
)

type CmdJobDto struct {
	Cmds []string `json:"cmds"`
	Envs []string `json:"envs"`
}

type AliveRequest struct {
	HostStatus   HostStatus   `json:"hostStatus"`
	RunnerStatus RunnerStatus `json:"runnerStatus"`
	JobsStatus   []JobsStatus `json:"jobsStatus"`
	AcceptJobs   []string     `json:"acceptJobs"`
	DenyJobs     []string     `json:"denyJobs"`
}

type RunnerStatus struct {
	RunningNum int `json:"runningNum"`
}

type HostStatus struct {
	// TODO CPU/Memory/Disk info
}

type JobsStatus struct {
	JobRunningId   string            `json:"jobRunningId"`
	AtomLogs       []AtomLog         `json:"atomLogs"`
	Finished       bool              `json:"finished"`
	FinishedStatus JobFinishedStatus `json:"finishedStatus"`
}

type JobFinishedStatus string

const (
	SUCCESS = JobFinishedStatus("SUCCESS")
	FAILURE = JobFinishedStatus("FAILURE")
)

// AtomLog 执行日志
type AtomLog struct {
	LogType   int    `json:"logType"`
	LogBody   string `json:"logBody"`
	OrderId   int    `json:"jobOrderId"`
	TimeStamp int64  `json:"timestamp"`
}

func (api *RunnerHttpApi) AliveToServer(aliveRequest *AliveRequest) (*AliveResponse, error) {
	var response AliveResponse
	err := api.doHttp("POST", api.address+"/live", aliveRequest, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
