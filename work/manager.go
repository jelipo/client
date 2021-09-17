package work

import (
	"client/util"
	"errors"
	"log"
	"sync"
)

type JobManager struct {
	status              ManagerStatus
	statusLock          sync.RWMutex
	runningMap          map[string]WorkerRunningStatus
	acceptRunningJobIds []string
	denyRunningJobIds   []string
}

type ManagerStatus struct {
	runningNum int
}

type WorkerRunningStatus struct {
	flag    util.AsyncRunFlag
	starter *JobStarter
}

type JobOutStatus struct {
	JobRunningId string
	AtomLogs     []AtomLog
	Finished     bool
	Success      bool
}

func NewWorkerManager(maxNum int) JobManager {
	return JobManager{
		status:              ManagerStatus{runningNum: 0},
		statusLock:          sync.RWMutex{},
		runningMap:          make(map[string]WorkerRunningStatus, 128),
		acceptRunningJobIds: []string{},
		denyRunningJobIds:   []string{},
	}
}

func (manager *JobManager) ReadStatus() (ManagerStatus, map[string]JobOutStatus) {
	manager.statusLock.RLock()
	defer manager.statusLock.RUnlock()
	status := manager.status
	var statusMap = map[string]JobOutStatus{}
	for JobRunningId, workerStatus := range manager.runningMap {
		stepLog := workerStatus.starter.JobLog()
		logs := stepLog.GetLogs(100)
		outStatus := JobOutStatus{
			JobRunningId: JobRunningId,
			AtomLogs:     logs,
			Finished:     workerStatus.flag.IsDone(),
			Success:      !workerStatus.flag.HaveError(),
		}
		statusMap[JobRunningId] = outStatus
	}
	// 更新状态
	for _, outStatus := range statusMap {
		if outStatus.Finished {
			delete(manager.runningMap, outStatus.JobRunningId)
			manager.status.runningNum = manager.status.runningNum - 1
		}
	}
	return status, statusMap
}

func (manager *JobManager) AddAcceptDenyRunningJobId(acceptIds []string, denyIds []string) {
	manager.denyRunningJobIds = append(manager.denyRunningJobIds, denyIds...)
	manager.acceptRunningJobIds = append(manager.acceptRunningJobIds, acceptIds...)
}

func (manager *JobManager) GetAndCleanAcceptDenyRunningJobIds() ([]string, []string) {
	accept := manager.acceptRunningJobIds
	deny := manager.denyRunningJobIds
	manager.denyRunningJobIds = []string{}
	manager.acceptRunningJobIds = []string{}
	return accept, deny
}

func (manager *JobManager) AddNewJob(jobRunningId string, sources []Source, newJob *NewJob) error {
	manager.statusLock.Lock()
	defer manager.statusLock.Unlock()
	// check status
	_, ok := manager.runningMap[jobRunningId]
	if ok {
		return errors.New("'" + jobRunningId + "' already exited")
	}
	// creat new work
	starter, err := NewJobStarter(sources, newJob)
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

func asyncRunWorker(starter *JobStarter) util.AsyncRunFlag {
	return util.NewAsyncRunFlag(func() error {
		defer clean(starter)
		err := run(starter)
		if err != nil {
			return err
		}
		log.Println("Job:" + starter.jobRunningId + " finished.")
		return nil
	})
}

func clean(starter *JobStarter) {
	err := starter.pipeJobDir.CleanRunningJobDir()
	if err != nil {
		log.Println("Clean job dir failed jobRunningId: " + starter.jobRunningId + " err:" + err.Error())
		return
	} else {
		log.Println("Clean job dir success jobRunningId:" + starter.jobRunningId)
	}
}

func run(starter *JobStarter) error {
	err := starter.RunStarter()
	if err != nil {
		log.Println("Run job:" + starter.jobRunningId + " failed. " + err.Error())
		return err
	}
	return nil
}
