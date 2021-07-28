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
	flag    util.GoroutinesFlag
	starter *WorkerStarter
}

type WorkerOutStatus struct {
	pipeSoleId string
	atomLogs   []AtomLog
	done       bool
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
	for pipeSoleId, workerStatus := range manager.runningMap {
		stepLog := workerStatus.starter.StepLog()
		logs := stepLog.GetLogs(100)
		outStatus := WorkerOutStatus{
			pipeSoleId: pipeSoleId,
			atomLogs:   logs,
			done:       workerStatus.flag.IsDone(),
		}
		statusMap[pipeSoleId] = outStatus
	}
	for _, outStatus := range statusMap {
		if outStatus.done {
			delete(manager.runningMap, outStatus.pipeSoleId)
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
	_, ok := manager.runningMap[newWork.PipeSoleId]
	if ok {
		return errors.New("'" + newWork.PipeSoleId + "' already exited")
	}
	// creat new work
	starter, err := NewWorkerStarter(sources, newWork)
	if err != nil {
		return err
	}
	flag := goRunWorker(starter)
	manager.runningMap[newWork.PipeSoleId] = WorkerStatus{
		flag:    flag,
		starter: starter,
	}
	manager.status.runningNum += 1
	return nil
}

func goRunWorker(starter *WorkerStarter) util.GoroutinesFlag {
	return util.NewGoroutinesFlagAndRun(func() {
		run(starter)
	})
}

func run(starter *WorkerStarter) {
	err := starter.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
