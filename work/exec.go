package work

import (
	"client/util"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

//go:embed exec.bash
var baseBashBytes []byte

type Executor struct {
	workDir   string
	actionLog *ActionLog
	customEnv []string
}

func NewExec(workDir string, actionLog *ActionLog, customEnv []string) Executor {
	return Executor{workDir: workDir, actionLog: actionLog, customEnv: customEnv}
}

func (executor *Executor) ExecShell(shellPart string) error {
	randomStr := util.RandNumAndLettersStr(10)
	bashFilePath := executor.workDir + "/" + randomStr + ".bash"
	envFileName := randomStr + ".env"
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
		listenLog(stdOutPipe, executor.actionLog)
	})

	goPackage.AddAndRun(func() {
		listenLog(stdErrPipe, executor.actionLog)
	})
	goPackage.WaitAllDone()
	waitErr := cmd.Wait()
	if waitErr != nil {
		return err
	}
	code := cmd.ProcessState.ExitCode()
	if code != 0 {
		return errors.New(fmt.Sprintf("Error,shell exited code is %d", code))
	}
	return nil
}

func listenLog(stdPipe io.ReadCloser, actionLog *ActionLog) {
	buf := make([]byte, 4096)
	for true {
		time.Sleep(time.Duration(50) * time.Millisecond)
		read, err := stdPipe.Read(buf)
		if err != nil {
			fmt.Println("listen executor log failed")
			break
		}
		if read <= 0 {
			fmt.Println("listen executor finished")
			break
		}
		log := string(buf)
		actionLog.AddExecLog(log)
	}
}
