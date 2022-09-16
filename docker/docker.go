package docker

import (
	"fmt"

	"github.com/moby/moby/client"
)

func New() *docker {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err)
	}
	return &docker{
		cli: cli,
	}
}

type docker struct {
	cli *client.Client
}
