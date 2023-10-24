package docker

import (
	"github.com/arshamalh/dockeroller/log"
	"github.com/moby/moby/client"
)

func New() *docker {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	return &docker{
		cli: cli,
	}
}

type docker struct {
	cli *client.Client
}
