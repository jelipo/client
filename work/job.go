package work

type NewJob struct {
	//NewWork       work.NewWork  `json:"newWork"`
	JobRunningId  string    `json:"jobRunningId"`
	Sources       []Source  `json:"sources"`
	PipeId        string    `json:"pipeId"`
	StageId       string    `json:"stageId"`
	JobId         string    `json:"jobId"`
	PipeRunningId string    `json:"pipeRunningId"`
	MainSourceId  string    `json:"mainSourceId"`
	JobType       string    `json:"jobType"`
	CmdJobDto     CmdJobDto `json:"cmdJob"` // Exited when jobType is "COMMAND"
}

type Source struct {
	SourceId        string          `json:"sourceId"`
	SourceType      string          `json:"sourceType"`
	UseCache        bool            `json:"useCache"`
	ProjectName     string          `json:"projectName"`
	GitSourceConfig GitSourceConfig `json:"gitSourceConfig"`
}

type GitSourceConfig struct {
	GitAddress        string `json:"gitAddress"`
	AuthType          string `json:"authType"` //拉取git的身份验证方式,PublicKey or Password
	Branch            string `json:"branch"`
	CommitId          string `json:"commitId"` //CommitHashId
	AuthUsername      string `json:"authUsername"`
	AuthPassword      string `json:"authPassword"`
	AuthPublicKeyStr  string `json:"authPublicKeyStr"`
	AuthPublicKeyPath string `json:"authPublicKeyPath"`
}

type CmdJobDto struct {
	Cmds []string `json:"cmds"`
	Envs []string `json:"envs"`
}
