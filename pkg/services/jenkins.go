package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bndr/gojenkins"
	"time"
)

type Jenkins struct {
	Server   *gojenkins.Jenkins
	Username string `json:"username"` // admin
	Token    string `json:"token"`    // 11dc2e1df59b44ac920d420e40c08e3e56
	Host     string `json:"host"`     // "http://localhost:8080/"
}

func NewJenkins(username, token, host string) (Jenkins, error) {
	var j = Jenkins{Username: username, Token: token, Host: host}
	err := j.GetJenkins()
	if err != nil {
		return j, err
	}
	return j, nil
}

func (j *Jenkins) Build(jenkins *gojenkins.Jenkins, jobName string, options map[string]string) error {
	queueId, err := jenkins.BuildJob(context.TODO(), jobName, options)
	fmt.Println("queueId:", queueId)
	if err != nil {
		fmt.Println("BuildJob:", err)
		return err
	}
	var task *gojenkins.Task
	var flag bool
	for !flag {
		time.Sleep(2 * time.Second)
		task, err = jenkins.GetQueueItem(context.TODO(), queueId)
		if err != nil {
			fmt.Println("GetQueueItem:", err)
			return err
		}
		//fmt.Println("Executable.Number:", task.Raw.Executable.Number)
		fmt.Println("Pending:", task.Raw.Pending, task.Raw.Why)
		//tr, _ := json.Marshal(task.Raw)
		//fmt.Println("task raw:", string(tr))
		if task.Raw.Executable.Number != 0 {
			flag = true
		}
		//fmt.Println("task raw:", i)

	}
	fmt.Println("GetQueueItem finish! build number:", task.Raw.Executable.Number)

	var build *gojenkins.Build
	flag = false
	for !flag {
		time.Sleep(2 * time.Second)
		build, err = jenkins.GetBuild(context.TODO(), jobName, task.Raw.Executable.Number)
		if err != nil {
			fmt.Println("GetBuild:", err)
			return err
		}
		fmt.Println("GetBuild: ", build.Info().EstimatedDuration)
		if len(build.Raw.Result) != 0 && !build.Raw.Building {
			flag = true
		}
	}
	br, _ := json.Marshal(build.Raw)
	fmt.Println("build raw:", build.Raw.Result)
	fmt.Println("build raw:", build.Raw.Duration)
	fmt.Println("build raw:", build.Raw.Number)
	fmt.Println("build raw:", string(br))
	//fmt.Println("GetBuild finish! ", build.Info())
	//fmt.Println("GetBuild finish! ", build.Info().Result)

	return nil
}

func (j *Jenkins) GetJenkins() error {
	var err error
	j.Server, err = gojenkins.CreateJenkins(nil, j.Host, j.Username, j.Token).Init(context.TODO())
	if err != nil {
		return err
	}
	return nil
}
