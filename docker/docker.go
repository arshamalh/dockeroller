package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Docker is everything dockeroller needs from the Docker daemon. Handlers
// depend on this interface so tests can swap in a fake.
type Docker interface {
	GetContainer(ctx context.Context, containerID string) (*entities.Container, error)
	ContainersList(ctx context.Context, filters filters.Args) ([]*entities.Container, error)
	// ContainerLogs follows the container's stdout and stderr as plain text
	// (already demultiplexed). Cancel ctx to stop following.
	ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error)
	// ContainerStats streams one JSON stats document per line until ctx is
	// cancelled.
	ContainerStats(ctx context.Context, containerID string) (io.ReadCloser, error)
	ContainerStart(ctx context.Context, containerID string) error
	ContainerStop(ctx context.Context, containerID string) error
	ContainerRemove(ctx context.Context, containerID string, removeForm *entities.ContainerRemoveForm) error
	ContainerRename(ctx context.Context, containerID, newName string) error
	ContainerPause(ctx context.Context, containerID string) error
	ContainerUnpause(ctx context.Context, containerID string) error
	ImagesList(ctx context.Context) ([]*entities.Image, error)
	ImageTag(ctx context.Context, imageID, newTag string) error
	ImageRemove(ctx context.Context, imageID string, force, pruneChildren bool) error
}

// New connects to the daemon configured by the standard environment
// (DOCKER_HOST, DOCKER_TLS_VERIFY, DOCKER_CERT_PATH) and negotiates the API
// version, so it works with older and newer daemons alike.
func New() (*docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}
	return &docker{cli: cli}, nil
}

// Ping checks that the daemon is reachable.
func (d *docker) Ping(ctx context.Context) error {
	_, err := d.cli.Ping(ctx)
	return err
}

// Close releases the client's connections.
func (d *docker) Close() error {
	return d.cli.Close()
}

type docker struct {
	cli *client.Client
}
