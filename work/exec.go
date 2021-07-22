package work

import (
	"client/util"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

//go:embed exec.bash
var baseBashBytes []byte

type Executor struct {
	workDir   string
	actionLog *ActionLog
	customEnv []string
	// mill_second
	timeOut  int
	stopFlag bool
}

func NewExec(workDir string, actionLog *ActionLog, customEnv []string, timeOut int) Executor {
	return Executor{
		workDir:   workDir,
		actionLog: actionLog,
		customEnv: customEnv,
		timeOut:   timeOut,
		stopFlag:  false,
	}
}

func (executor *Executor) ExecShell(shellPart string) error {
	randomStr := util.RandLowcaseLetters(10)
	bashFilePath := executor.workDir + "/" + randomStr + ".bash"
	envFileName := randomStr + ".env"
	defer os.Remove(bashFilePath)
	defer os.Remove(executor.workDir + "/" + envFileName)
	bashFile, err := os.Create(bashFilePath)
	if err != nil {
		return err
	}
	_, err = bashFile.Write(baseBashBytes)
	if err != nil {
		return err
	}
	cmd := exec.Command("bash", bashFilePath, shellPart, envFileName)
	cmd.Dir = executor.workDir
	// Add custom env
	env := cmd.Env
	env = append(env, executor.customEnv...)
	cmd.Env = env
	// log pipe
	stdErrPipe, errErr := cmd.StderrPipe()
	if errErr != nil {
		return errErr
	}
	stdOutPipe, outErr := cmd.StdoutPipe()
	if outErr != nil {
		return outErr
	}
	goPackage := util.NewGoPackage()
	goPackage.AddAndRun(func() {
		listenLog(&stdOutPipe, executor.actionLog)
	})

	goPackage.AddAndRun(func() {
		listenLog(&stdErrPipe, executor.actionLog)
	})
	startErr := cmd.Start()
	if startErr != nil {
		return startErr
	}
	goPackage.WaitAllDone()
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
