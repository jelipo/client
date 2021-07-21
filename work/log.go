package work

// ExecLog 执行日志
type ExecLog struct {
	LogType int
	// 日志实体
	LogBody string
}

const (
	ActionLogType = 1
	StepLogType   = 2
	SysLogType    = 3
)

type StepLog struct {
	logChannel chan ExecLog
}

func NewStepLog() StepLog {
	return StepLog{logChannel: make(chan ExecLog, 1024)}
}

func (stepLog *StepLog) NewAction(stepName string) ActionLog {
	stepLog.logChannel <- ExecLog{LogType: StepLogType, LogBody: stepName}
	return ActionLog{StepLogChannel: &stepLog.logChannel}
}

type ActionLog struct {
	StepLogChannel *chan ExecLog
}

func (actionLog ActionLog) AddExecLog(logBody string) {
	*actionLog.StepLogChannel <- ExecLog{LogType: ActionLogType, LogBody: logBody}
}

func (actionLog ActionLog) AddSysLog(logBody string) {
	*actionLog.StepLogChannel <- ExecLog{LogType: SysLogType, LogBody: logBody}
}
