package work

import (
	"client/api"
	"log"
	"time"
)

const (
	ActionLogType  = 1
	ActionNameType = 2
	SysLogType     = 3
)

type JobLog struct {
	logChannel chan api.AtomLog
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
		logChannel: make(chan api.AtomLog, 1024),
		orderGen:   OrderIdGen{orderIdTemp: 0},
	}
	return jobLog
}

func (jobLog *JobLog) NewAction(actionName string) ActionLog {
	jobLog.logChannel <- api.AtomLog{
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

func (jobLog *JobLog) GetLogs(maxSize int) []api.AtomLog {
	chanLen := len(jobLog.logChannel)
	if chanLen == 0 {
		return make([]api.AtomLog, 0)
	}
	var buffer []api.AtomLog
	if chanLen < maxSize {
		buffer = make([]api.AtomLog, chanLen)
	} else {
		buffer = make([]api.AtomLog, maxSize)
	}
	size := len(buffer)
	for i := 0; i < size; i++ {
		log := <-jobLog.logChannel
		buffer[i] = log
	}
	return buffer
}

type ActionLog struct {
	StepLogChannel *chan api.AtomLog
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
	*actionLog.StepLogChannel <- api.AtomLog{
		LogType:   logType,
		LogBody:   logBody,
		OrderId:   actionLog.orderIdGen.GetAndAdd(),
		TimeStamp: time.Now().UnixMilli(),
	}
}
