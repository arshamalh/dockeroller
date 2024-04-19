package docker

import (
	"context"
	"io"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/docker/docker/api/types/filters"
	"github.com/moby/moby/client"
)

type Docker interface {
	// TODO: Add context to these functions, as they all accept context in their lower level
	GetContainer(ctx context.Context, containerID string) (*entities.Container, error)
	ContainersList(ctx context.Context, filters filters.Args) ([]*entities.Container, error)
	// TODO: context for these one is so important as we might be able to remove quite channel by its help
	ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error)
	ContainerStats(ctx context.Context, containerID string) (io.ReadCloser, error)
	ContainerStart(ctx context.Context, containerID string) error
	ContainerStop(ctx context.Context, containerID string) error
	ContainerRemove(ctx context.Context, containerID string, removeForm *entities.ContainerRemoveForm) error
	ContainerRename(ctx context.Context, containerID, newName string) error

	ImagesList(ctx context.Context, filters filters.Args) ([]*entities.Image, error)
	ImageTag(ctx context.Context, imageID, newTag string) error
	ImageRemove(ctx context.Context, imageID string, force, pruneChildren bool) error
}

func New() *docker {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
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
