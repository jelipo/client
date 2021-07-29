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
	runningMap map[string]WorkerStatus
}

type ManagerStatus struct {
	maxNum     int
	runningNum int
}

type WorkerStatus struct {
	flag    util.AsyncRunFlag
	starter *WorkerStarter
}

type WorkerOutStatus struct {
	workerId string
	atomLogs []AtomLog
	done     bool
}

func NewWorkerManager(maxNum int) Manager {
	return Manager{
		status:     ManagerStatus{},
		statusLock: sync.RWMutex{},
		runningMap: make(map[string]WorkerStatus, 128),
	}
}

func (manager *Manager) ReadStatus() ManagerStatus {
	manager.statusLock.RLock()
	defer manager.statusLock.RUnlock()
	status := manager.status
	var statusMap = make(map[string]WorkerOutStatus, status.maxNum)
	for workerId, workerStatus := range manager.runningMap {
		stepLog := workerStatus.starter.StepLog()
		logs := stepLog.GetLogs(100)
		outStatus := WorkerOutStatus{
			workerId: workerId,
			atomLogs: logs,
			done:     workerStatus.flag.IsDone(),
		}
		statusMap[workerId] = outStatus
	}
	for _, outStatus := range statusMap {
		if outStatus.done {
			delete(manager.runningMap, outStatus.workerId)
		}
	}
	return status
}

func (manager *Manager) AddNewWork(sources []Source, newWork *NewWork) error {
	manager.statusLock.Lock()
	defer manager.statusLock.Unlock()
	// check status
	status := &manager.status
	if status.runningNum >= status.maxNum {
		return errors.New("the task has reached the maximum limit")
	}
	_, ok := manager.runningMap[newWork.WorkerId]
	if ok {
		return errors.New("'" + newWork.WorkerId + "' already exited")
	}
	// creat new work
	starter, err := NewWorkerStarter(sources, newWork)
	if err != nil {
		return err
	}
	flag := asyncRunWorker(starter)
	manager.runningMap[newWork.WorkerId] = WorkerStatus{
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
