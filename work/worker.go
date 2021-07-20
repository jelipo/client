package work

type NewWork struct {
	Id   string
	Type int
}

const (
	CommandType = 1
	DeployType  = 2
)

type CommandWorkType struct {
	Sources  []Source
	Cmds     []string
	WorkPath string
}

type Worker struct {
	Num int
}

func NewWorker(newWork *NewWork) {

}
