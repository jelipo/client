package work

import (
	"client/util"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Executor struct {
	workDir    string
	actionLog  *ActionLog
	appendEnvs []string
	recordEnv  bool
	// mill_second
	timeOut  int
	stopFlag bool
	stopChan *chan bool
}

func NewExec(workDir string, actionLog *ActionLog, customEnv []string, timeOut int, recordEnv bool, stopChan *chan bool) Executor {
	return Executor{
		workDir:    workDir,
		actionLog:  actionLog,
		appendEnvs: customEnv,
		recordEnv:  recordEnv,
		timeOut:    timeOut,
		stopFlag:   false,
		stopChan:   stopChan,
	}
}

func (executor *Executor) ExecShell(shellPart string) error {
	randomStr := util.RandLowcaseLetters(10)
	envAbFilePath := executor.workDir + "/" + randomStr + ".env"
	defer os.Remove(envAbFilePath)
	//cmd := exec.Command("bash", bashFilePath, shellPart, envFileName)
	var fullCmd = shellPart
	if executor.recordEnv {
		fullCmd = shellPart + " env >> \"" + envAbFilePath + "\""
	}
	cmd := exec.Command("bash", "-c", "source /etc/profile && "+fullCmd)
	cmd.Dir = executor.workDir
	// Add custom env
	env := cmd.Env
	env = append(env, executor.appendEnvs...)
	cmd.Env = env
	// log pipe
	startErr := executor.startAndWait(cmd)
	if startErr != nil {
		return startErr
	}
	if executor.recordEnv {
		// Read environments
		afterEnvs, envErr := readEnvs(envAbFilePath)
		if envErr != nil {
			return envErr
		}
		executor.appendEnvs = afterEnvs
	}
	return nil
}

func (executor *Executor) startAndWait(cmd *exec.Cmd) error {
	stdErrPipe, errErr := cmd.StderrPipe()
	if errErr != nil {
		return errErr
	}
	stdOutPipe, outErr := cmd.StdoutPipe()
	if outErr != nil {
		return outErr
	}
	goPackage := util.NewAsyncCollect()
	goPackage.AddAndRun(func() {
		listenLog(&stdOutPipe, executor.actionLog)
	})
	goPackage.AddAndRun(func() {
		listenLog(&stdErrPipe, executor.actionLog)
	})
	if startErr := cmd.Start(); startErr != nil {
		return startErr
	}
	for true {
		if goPackage.IsAllDone() {
			break
		} else {
			if len(*executor.stopChan) != 0 {
				if err := cmd.Process.Kill(); err != nil {
					return err
				}
			} else {
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
		}
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		return waitErr
	}
	code := cmd.ProcessState.ExitCode()
	if code != 0 {
		return errors.New(fmt.Sprintf("Error,shell exited code is %d", code))
	}
	return nil
}

func readEnvs(envAbFilePath string) ([]string, error) {
	envBytes, readErr := os.ReadFile(envAbFilePath)
	if readErr != nil {
		return nil, readErr
	}
	envsStr := string(envBytes)
	envs := strings.Split(envsStr, "\n")
	return envs, nil
}

func (executor *Executor) GetEnvs() *[]string {
	return &executor.appendEnvs
}

func listenLog(stdPipe *io.ReadCloser, actionLog *ActionLog) {
	buf := make([]byte, 1024)
	for true {
		read, err := (*stdPipe).Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("listen executor log failed" + err.Error())
			break
		}
		if read <= 0 {
			fmt.Println("listen executor finished")
			break
		}
		log := string(buf[:read])
		actionLog.AddExecLog(log)
	}
}
