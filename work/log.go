package work

// AtomLog 执行日志
type AtomLog struct {
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
	logChannel chan AtomLog
}

func NewStepLog() StepLog {
	return StepLog{logChannel: make(chan AtomLog, 1024)}
}

func (stepLog *StepLog) NewAction(stepName string) ActionLog {
	stepLog.logChannel <- AtomLog{LogType: StepLogType, LogBody: stepName}
	return ActionLog{StepLogChannel: &stepLog.logChannel}
}

func (stepLog *StepLog) GetLogs(maxSize int) []AtomLog {
	chanLen := len(stepLog.logChannel)
	if chanLen == 0 {
		return nil
	}
	var buffer []AtomLog
	if chanLen < maxSize {
		buffer = make([]AtomLog, chanLen)
	} else {
		buffer = make([]AtomLog, maxSize)
	}
	size := len(buffer)
	for i := 0; i < size; i++ {
		log := <-stepLog.logChannel
		buffer[i] = log
	}
	return buffer
}

type ActionLog struct {
	StepLogChannel *chan AtomLog
}

func (actionLog ActionLog) AddExecLog(logBody string) {
	*actionLog.StepLogChannel <- AtomLog{LogType: ActionLogType, LogBody: logBody}
}

func (actionLog ActionLog) AddSysLog(logBody string) {
	*actionLog.StepLogChannel <- AtomLog{LogType: SysLogType, LogBody: logBody}
}

func (actionLog ActionLog) Write(bytes []byte) (n int, err error) {
	*actionLog.StepLogChannel <- AtomLog{LogType: ActionLogType, LogBody: string(bytes)}
	return len(bytes), nil
}
