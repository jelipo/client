package work

import (
	"log"
	"time"
)

// AtomLog 执行日志
type AtomLog struct {
	LogType int `json:"logType"`
	// 日志实体
	LogBody string `json:"logBody"`

	OrderId int `json:"jobOrderId"`

	TimeStamp int64 `json:"timestamp"`
}

const (
	ActionLogType  = 1
	ActionNameType = 2
	SysLogType     = 3
)

type JobLog struct {
	logChannel chan AtomLog
	orderGen   OrderIdGen
}

type OrderIdGen struct {
	orderIdTemp int
}

func (orderIdGen *OrderIdGen) GetAndAdd() int {
	orderIdGen.orderIdTemp = orderIdGen.orderIdTemp + 1
	return orderIdGen.orderIdTemp
}

func NewJobLog() JobLog {
	jobLog := JobLog{
		logChannel: make(chan AtomLog, 1024),
		orderGen:   OrderIdGen{orderIdTemp: 0},
	}
	return jobLog
}

func (jobLog *JobLog) NewAction(actionName string) ActionLog {
	jobLog.logChannel <- AtomLog{
		LogType:   ActionNameType,
		LogBody:   actionName + "\n",
		TimeStamp: time.Now().UnixMilli(),
		OrderId:   jobLog.orderGen.GetAndAdd(),
	}
	return ActionLog{
		StepLogChannel: &jobLog.logChannel,
		orderIdGen:     &jobLog.orderGen,
	}
}

func (jobLog *JobLog) GetLogs(maxSize int) []AtomLog {
	chanLen := len(jobLog.logChannel)
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
		log := <-jobLog.logChannel
		buffer[i] = log
	}
	return buffer
}

type ActionLog struct {
	StepLogChannel *chan AtomLog
	orderIdGen     *OrderIdGen
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
	log.Print(logBody)
	*actionLog.StepLogChannel <- AtomLog{
		LogType:   logType,
		LogBody:   logBody,
		OrderId:   actionLog.orderIdGen.GetAndAdd(),
		TimeStamp: time.Now().UnixMilli(),
	}
}
