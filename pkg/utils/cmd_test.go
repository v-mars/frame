package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

//Golang 执行 shell脚本，接收返回值
func TestAA(t *testing.T) {
	//Golang 执行 shell 脚本，不接收返回值
	// 返回一个 cmd 对象
	cmd := exec.Command("sh", "-c", "./scripts/curl.sh")

	// 如果只执行命令，不接收返回值
	err := cmd.Run()
	if err != nil {
		return
	}

}
func TestA1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", "nginx -V")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		fmt.Println(stderr.String())

	}
	fmt.Println("Result: " + out.String())
	fmt.Println("stderr: " + stderr.String())
	fmt.Println("state", cmd.ProcessState)
	fmt.Println("state ExitCode", cmd.ProcessState.ExitCode())
	fmt.Println("state String", cmd.ProcessState.String())

}
func TestA2(t *testing.T) {
	out, err := exec.Command("/usr/local/bin/docker", "pss").Output()
	if err != nil {
		//fmt.Println(err)
		log.Fatal(err)
	}

	fmt.Printf("The date is %s\n", out)

}

// Golang 执行 shell脚本，接收返回值
func TestB(t *testing.T) {
	// 返回一个 cmd 对象
	cmd := exec.Command("docker", "pss")

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	//stdOut = cmd.Stdout
	//cmd.Stderr = &stdErr

	// 收返回值[]byte, error
	b, err := cmd.Output()
	if err != nil {
		fmt.Println("err:", err)
	}

	fmt.Println(string(b))
	fmt.Println(stdOut.String())
	fmt.Println(stdErr.String())

}

func TestMkdir(t *testing.T) {
	//getwd, err := os.Getwd()
	//if err != nil {
	//	fmt.Println("getwd err:", err)
	//	return
	//}
	//fmt.Println("getwd:", getwd)
	getwd := "/Users/xxx/work/dev/go/air"
	dir := getwd + "/data/program/goapp/golang"
	err := os.Mkdir(dir, 0755)
	if err != nil {
		fmt.Println(err)
	}

}

func TestMkdirAll(t *testing.T) {
	getwd := "/Users/xxx/work/dev/go/air"
	dir := getwd //+ "/data/program/goapp/golang"
	_, err := os.Stat(dir)
	if err != nil {
		fmt.Println("err:", err)
		fmt.Println("os.IsNotExist(err):", os.IsNotExist(err))
	}

	//err := os.MkdirAll(dir, 0755)

	//if err != nil {
	//
	//	fmt.Println(err)
	//
	//}

}

func TestDirRemove(t *testing.T) {
	getwd := "/Users/xxx/work/dev/go/air"
	dir := getwd + "/data/program/goapp/golang"
	err := os.Remove(dir)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

}

func TestTime(t *testing.T) {
	duration, err := time.ParseDuration("5d3m")
	if err != nil {
		fmt.Println("duration err:", err)
		return
	}
	fmt.Println("duration:", duration)

}

// Golang 执行 shell脚本，并实时打印 shell 脚本输出日志信息
func TestC(t *testing.T) {
	execute()

}

func asyncLog(reader io.ReadCloser) error {
	cache := ""
	buf := make([]byte, 1024, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed") {
				err = nil
			}
			return err
		}
		if num > 0 {
			oByte := buf[:num]
			//h.logInfo = append(h.logInfo, oByte...)
			oSlice := strings.Split(string(oByte), "\n")
			line := strings.Join(oSlice[:len(oSlice)-1], "\n")
			fmt.Printf("%s%s\n", cache, line)
			cache = oSlice[len(oSlice)-1]
		}
	}
	return nil
}

func execute() error {
	cmd := exec.Command("sh", "-c", "./scripts/curl.sh")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %s......", err.Error())
		return err
	}

	go asyncLog(stdout)
	go asyncLog(stderr)

	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for command execution: %s......", err.Error())
		return err
	}

	return nil
}
