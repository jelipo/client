package work

type NewJob struct {
	//NewWork       work.NewWork  `json:"newWork"`
	JobRunningId  string   `json:"jobRunningId"`
	Sources       []Source `json:"sources"`
	PipeId        string   `json:"pipeId"`
	StageId       string   `json:"stageId"`
	JobId         string   `json:"jobId"`
	PipeRunningId string   `json:"pipeRunningId"`
	MainSourceId  string   `json:"mainSourceId"`
	JobType       string   `json:"jobType"`
	CmdJobDto     string   `json:"cmdJob"` // Exited when jobType is "COMMAND"
}
