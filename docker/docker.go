package docker

import (
	"io"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/models"
	"github.com/moby/moby/client"
)

type Docker interface {
	GetContainer(containerID string) (*models.Container, error)
	ContainersList() []*models.Container
	ContainerLogs(containerID string) (io.ReadCloser, error)
	ContainerStats(containerID string) (io.ReadCloser, error)
	ContainerStart(containerID string) error
	ContainerStop(containerID string) error
	ContainerRemove(containerID string, removeForm *models.ContainerRemoveForm) error
	ContainerRename(containerID, newName string) error

	ImagesList() []*models.Image
}

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
