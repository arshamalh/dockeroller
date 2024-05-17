package docker

import (
	"context"
	"io"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func (d *docker) ContainersList(ctx context.Context, filters filters.Args) ([]*entities.Container, error) {
	raw_containers, err := d.cli.ContainerList(ctx,
		container.ListOptions{
			All:     true,
			Filters: filters,
		},
	)
	if err != nil {
		return nil, err
	}

	containers := make([]*entities.Container, 0)
	for _, raw_cont := range raw_containers {
		containers = append(containers, &entities.Container{
			ID:     raw_cont.ID,
			Name:   raw_cont.Names[0],
			Image:  raw_cont.Image,
			Status: raw_cont.Status,
			State:  entities.ContainerState(raw_cont.State),
		})
	}
	return containers, nil
}

func (d *docker) GetContainer(ctx context.Context, containerID string) (*entities.Container, error) {
	container, err := d.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return &entities.Container{
		ID:     container.ID,
		Name:   container.Name,
		Image:  container.Image,
		Status: container.State.Status,
		State:  entities.ContainerState(container.State.Status),
	}, nil
}

func (d *docker) ContainerStats(ctx context.Context, containerID string) (io.ReadCloser, error) {
	stats, err := d.cli.ContainerStats(ctx, containerID, true)
	return stats.Body, err
}

func (d *docker) ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	// TODO: Interesting options about logs are available, you can get them from user settings
	return d.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Details:    false,
	})
}

func (d *docker) ContainerStart(ctx context.Context, containerID string) error {
	return d.cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (d *docker) ContainerStop(ctx context.Context, containerID string) error {
	return d.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (d *docker) ContainerRemove(ctx context.Context, containerID string, removeForm *entities.ContainerRemoveForm) error {
	return d.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: removeForm.RemoveVolumes,
		Force:         removeForm.Force,
	})
}

func (d *docker) ContainerRename(ctx context.Context, containerID, newName string) error {
	return d.cli.ContainerRename(ctx, containerID, newName)
}
