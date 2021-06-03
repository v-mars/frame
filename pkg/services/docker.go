package services

import (
	"fmt"
	"github.com/docker/docker/client"
)

type Docker struct {
	Host    string
	Version string
	Client  *client.Client
}

func NewClient(host, port,version string) (*client.Client, error) {
	Host := fmt.Sprintf("tcp://%s:%s", host, port)
	clientAPi,err := client.NewClient(Host, version,nil, nil )
	if err != nil {
		return nil, err
	}
	return clientAPi, nil
}

func NewClientDocker(host, port,version string) (*client.Client, error) {
	Host := fmt.Sprintf("tcp://%s:%s", host, port)

	clientAPi,err := client.NewClient(Host, version,nil, nil )
	if err != nil {
		return nil, err
	}
	return clientAPi, nil
}