package work

// AtomLog 执行日志
type AtomLog struct {
	LogType int `json:"logType"`
	// 日志实体
	LogBody string `json:"logBody"`

	OrderId int `json:"orderId"`
}

const (
	ActionLogType  = 1
	ActionNameType = 2
	SysLogType     = 3
)

type JobLog struct {
	logChannel chan AtomLog
}

func NewStepLog() JobLog {
	return JobLog{logChannel: make(chan AtomLog, 1024)}
}

func (stepLog *JobLog) NewAction(actionName string) ActionLog {
	stepLog.logChannel <- AtomLog{LogType: ActionNameType, LogBody: actionName + "\n"}
	return ActionLog{StepLogChannel: &stepLog.logChannel, orderIdTemp: 0}
}

func (stepLog *JobLog) GetLogs(maxSize int) []AtomLog {
	chanLen := len(stepLog.logChannel)
	if chanLen == 0 {
		return make([]AtomLog, 0)
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
	orderIdTemp    int
}

func (actionLog *ActionLog) AddExecLog(logBody string) {
	actionLog.internalAdd(ActionLogType, logBody)
}

func (actionLog *ActionLog) AddSysLog(logBody string) {
	actionLog.internalAdd(SysLogType, logBody)
}

func (actionLog *ActionLog) Write(bytes []byte) (n int, err error) {
	actionLog.internalAdd(ActionLogType, string(bytes))
	return len(bytes), nil
}

func (actionLog *ActionLog) internalAdd(logType int, logBody string) {
	actionLog.orderIdTemp = actionLog.orderIdTemp + 1
	*actionLog.StepLogChannel <- AtomLog{LogType: logType, LogBody: logBody, OrderId: actionLog.orderIdTemp}
}
