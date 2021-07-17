package main

import (
	"client/util"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {

	reader := newConfigReader("/home/cao/go/client/config.json")
	config, readErr := reader.readNewConfig()
	if readErr != nil {
		fmt.Println(readErr)
		return
	}
	fmt.Println(config)

	cmd := exec.Command("/bin/bash", "tem.sh")
	cmd.Dir = "/home/cao/go/client"
	stdErrPipe, errErr := cmd.StderrPipe()
	if errErr != nil {
		return
	}
	stdOutPipe, ouErr := cmd.StdoutPipe()
	if ouErr != nil {
		return
	}

	env := cmd.Env
	env = append(env, "TEST2=1")
	cmd.Env = env
	err := cmd.Start()
	if err != nil {
		println(err.Error())
		return
	}

	goPackage := util.NewGoPackage()
	goPackage.AddAndRun(func() {
		printLog(stdOutPipe)
	})

	goPackage.AddAndRun(func() {
		printLog(stdErrPipe)
	})
	goPackage.WaitAllDone()
	errErr = cmd.Wait()

	if cmd.ProcessState.Exited() {
		if cmd.ProcessState.ExitCode() != 0 {
			fmt.Println("Failed! Exit code is ", cmd.ProcessState.ExitCode())
		} else {
			fmt.Println("Success!")
		}
	}

	fmt.Println("All Done")
	println(os.Getenv("JAVA_HOME"))
	println(os.Getenv("TEST"))
	for i := range os.Environ() {
		println(os.Environ()[i])
	}
}

func printLog(stdPipe io.ReadCloser) {
	buf := make([]byte, 1024)
	for true {
		read, _ := stdPipe.Read(buf)
		if read <= 0 {
			break
		}
		fmt.Print(string(buf[:read]))
	}
}
