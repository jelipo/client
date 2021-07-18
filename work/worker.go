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
	Source   Source
	Cmds     []string
	WorkPath string
}

type GitCommandWorkType struct {
	GitRepo  string
	Cmds     []string
	WorkPath string
}

type Source struct {
	Type          int
	SourceAddress string
	SourceAttr    interface{}
}

type Worker struct {
	Num int
}
