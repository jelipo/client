package main

import (
	"fmt"
	"io"
	"os/exec"
)

func main() {
	cmd := exec.Command("mvn", "clean", "package", "-DskipTests", "-PDaoCloud")
	cmd.Dir = "E:\\idea\\dps-backend"
	stdErrPipe, _ := cmd.StderrPipe()
	stdOutPipe, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		println(err.Error())
		return
	}

	goPackage := NewGoPackage()
	goPackage.AddAndRun(func() {
		printLog(stdOutPipe)
	})

	goPackage.AddAndRun(func() {
		printLog(stdErrPipe)
	})
	goPackage.WaitAllDone()

	fmt.Println("All Done")
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
