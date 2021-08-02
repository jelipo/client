package work

import (
	"client/util"
	"errors"
	"fmt"
	"sync"
)

type Manager struct {
	status     ManagerStatus
	statusLock sync.RWMutex
	runningMap map[string]WorkerRunningStatus
}

type ManagerStatus struct {
	maxNum     int
	runningNum int
}

type WorkerRunningStatus struct {
	flag    util.AsyncRunFlag
	starter *WorkerStarter
}

type WorkerOutStatus struct {
	JobRunningId string
	AtomLogs     []AtomLog
	Done         bool
}

func NewWorkerManager(maxNum int) Manager {
	return Manager{
		status:     ManagerStatus{maxNum: maxNum, runningNum: 0},
		statusLock: sync.RWMutex{},
		runningMap: make(map[string]WorkerRunningStatus, 128),
	}
}

func (manager *Manager) ReadStatus() (ManagerStatus, map[string]WorkerOutStatus) {
	manager.statusLock.RLock()
	defer manager.statusLock.RUnlock()
	status := manager.status
	var statusMap = make(map[string]WorkerOutStatus, status.maxNum)
	for JobRunningId, workerStatus := range manager.runningMap {
		stepLog := workerStatus.starter.StepLog()
		logs := stepLog.GetLogs(100)
		outStatus := WorkerOutStatus{
			JobRunningId: JobRunningId,
			atomLogs:     logs,
			done:         workerStatus.flag.IsDone(),
		}
		statusMap[JobRunningId] = outStatus
	}
	for _, outStatus := range statusMap {
		if outStatus.done {
			delete(manager.runningMap, outStatus.JobRunningId)
		}
	}
	return status, statusMap
}

func (manager *Manager) AddNewJob(jobRunningId string, sources []Source, newWork *NewWork) error {
	manager.statusLock.Lock()
	defer manager.statusLock.Unlock()
	// check status
	status := &manager.status
	if status.runningNum >= status.maxNum {
		return errors.New("the task has reached the maximum limit")
	}
	_, ok := manager.runningMap[jobRunningId]
	if ok {
		return errors.New("'" + jobRunningId + "' already exited")
	}
	// creat new work
	starter, err := NewWorkerStarter(sources, newWork)
	if err != nil {
		return err
	}
	flag := asyncRunWorker(starter)
	manager.runningMap[jobRunningId] = WorkerRunningStatus{
		flag:    flag,
		starter: starter,
	}
	manager.status.runningNum += 1
	return nil
}

func asyncRunWorker(starter *WorkerStarter) util.AsyncRunFlag {
	return util.NewAsyncRunFlag(func() {
		run(starter)
	})
}

func run(starter *WorkerStarter) {
	err := starter.RunStarter()
	if err != nil {
		fmt.Println(err)
		return
	}
}
