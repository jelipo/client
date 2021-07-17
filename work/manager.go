package work

import (
	"errors"
	"sync"
)

type Manager struct {
	status     ManagerStatus
	statusLock sync.RWMutex
}

type ManagerStatus struct {
	maxNum     int
	runningNum int
}

func NewWorkerManager(maxNum int) Manager {
	return Manager{
		status: ManagerStatus{
			maxNum:     maxNum,
			runningNum: 0,
		},
	}
}

func (manager *Manager) readStatus() ManagerStatus {
	manager.statusLock.RLock()
	defer manager.statusLock.RUnlock()
	status := manager.status
	// TODO need deep copy
	return status
}

func (manager *Manager) addNewWork() error {
	manager.statusLock.Lock()
	defer manager.statusLock.Unlock()
	status := &manager.status
	if status.runningNum >= status.maxNum {
		return errors.New("已经到达了最高限制，不允许再添加")
	}
	// TODO add
	return nil
}
