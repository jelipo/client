package worker

import "sync"

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
