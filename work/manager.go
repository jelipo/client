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
	flag   *util.GoroutinesFlag
	worker *Worker
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
	// TODO need deep copy
	return status
}

func (manager *Manager) AddNewWork(newWork *NewWork) error {
	manager.statusLock.Lock()
	defer manager.statusLock.Unlock()
	// check status
	status := &manager.status
	if status.runningNum >= status.maxNum {
		return errors.New("the task has reached the maximum limit")
	}
	_, ok := manager.runningMap[newWork.StepId]
	if ok {
		return errors.New("'" + newWork.StepId + "' already exited")
	}
	// TODO creat new work

	manager.status.runningNum += 1
	return nil
}

func goRunWorker(newWork *NewWork) {
	util.NewGoroutinesFlagAndRun(func() {
		run(newWork)
	})
}

func run(newWork *NewWork) {
	worker, err := NewWorker(newWork)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = worker.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
