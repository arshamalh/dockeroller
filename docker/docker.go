package docker

import (
	"context"
	"io"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/moby/moby/client"
)

type Docker interface {
	GetContainer(containerID string) (*entities.Container, error)
	ContainersList() []*entities.Container
	ContainerLogs(containerID string) (io.ReadCloser, error)
	ContainerStats(containerID string) (io.ReadCloser, error)
	ContainerStart(containerID string) error
	ContainerStop(containerID string) error
	ContainerRemove(containerID string, removeForm *entities.ContainerRemoveForm) error
	ContainerRename(containerID, newName string) error

	ImagesList() []*entities.Image
	ImageTag(ctx context.Context, imageID, newTag string) error
	ImageRemove(ctx context.Context, imageID string, force, pruneChildren bool) error
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
