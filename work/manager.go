package work

import (
	"errors"
	"sync"
)

type Manager struct {
	status     ManagerStatus
	statusLock sync.RWMutex
	runningMap map[string]*Worker
}

type ManagerStatus struct {
	maxNum     int
	runningNum int
}

func NewWorkerManager(maxNum int) Manager {
	return Manager{
		runningMap: make(map[string]*Worker),
	}
}

func (manager *Manager) readStatus() ManagerStatus {
	manager.statusLock.RLock()
	defer manager.statusLock.RUnlock()
	status := manager.status
	// TODO need deep copy
	return status
}

func (manager *Manager) addNewWork(work *NewWork) error {
	manager.statusLock.Lock()
	defer manager.statusLock.Unlock()
	status := &manager.status
	if status.runningNum >= status.maxNum {
		return errors.New("the task has reached the maximum limit")
	}
	exitedWork := manager.runningMap[work.Id]
	if exitedWork != nil {
		return errors.New("'" + work.Id + "' already exited")
	}

	// TODO add
	manager.status.runningNum += 1
	return nil
}
