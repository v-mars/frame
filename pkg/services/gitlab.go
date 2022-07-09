package services

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type Gitlab struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Host     string `json:"host"`
	Client  *gitlab.Client `json:"client"`
	Project  *gitlab.Project `json:"project"`
}

func NewGitlab(host, token string) (Gitlab,error) {
	var g = Gitlab{Host: host,Token: token}
	err := g.GitlabClient()
	if err!=nil{
		return g, err
	}
	return g,nil
}

func (g *Gitlab) GitlabClient() error {
	urlSrt := fmt.Sprintf("%sapi/v4", g.Host)
	var err error

	g.Client, err = gitlab.NewClient(g.Token,gitlab.WithBaseURL(urlSrt), gitlab.WithoutRetries())

	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		return err
	}
	return nil
}

/*
@projectName: ops/java-demo
*/
func (g *Gitlab) GetProject(projectName string) error {
	tempSlice:= strings.Split(projectName, "/")
	if len(tempSlice) == 0{
		return errors.New("no project name")
	}
	name:=tempSlice[len(tempSlice)-1]

	pros, _, err := g.Client.Projects.ListProjects(&gitlab.ListProjectsOptions{Search: gitlab.String(name)})
	if err!=nil{
		return err
	}
	for _,v:=range pros {
		//fmt.Println("v:",v)
		if v.Name == name && v.PathWithNamespace == projectName {
			g.Project = v
			goto OutFor
		}

	}
OutFor:
	if len(g.Project.Name) == 0{
		return errors.New("not found project, plz check")
	}
	return nil
}

func (g *Gitlab) GetBranches() ([]*gitlab.Branch,error) {
	branches, _, err := g.Client.Branches.ListBranches(g.Project.ID,&gitlab.ListBranchesOptions{})
	if err!=nil{
		return branches,err
	}
	//fmt.Println("branches:", branches)
	return branches,nil
}

func (g *Gitlab) GetCommitIds() ([]*gitlab.Commit,error) {
	options := &gitlab.ListCommitsOptions{}
	commitIds, _, err := g.Client.Commits.ListCommits(g.Project.ID,options)
	if err!=nil{
		return commitIds,err
	}
	//fmt.Println("commitIds:", commitIds)
	return commitIds,nil
}