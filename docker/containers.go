package docker

import (
	"context"
	"io"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types"
)

func (d *docker) ContainersList() (containers []*entities.Container) {
	raw_containers, _ := d.cli.ContainerList(context.TODO(), types.ContainerListOptions{All: true})
	for _, raw_cont := range raw_containers {
		containers = append(containers, &entities.Container{
			ID:     raw_cont.ID,
			Name:   raw_cont.Names[0],
			Image:  raw_cont.Image,
			Status: raw_cont.Status,
			State:  entities.ContainerState(raw_cont.State),
		})
	}
	return
}

func (d *docker) GetContainer(containerID string) (*entities.Container, error) {
	container, err := d.cli.ContainerInspect(context.TODO(), containerID)
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

func (d *docker) ContainerStats(containerID string) (io.ReadCloser, error) {
	stats, err := d.cli.ContainerStats(context.TODO(), containerID, true)
	return stats.Body, err
}

func (d *docker) ContainerLogs(containerID string) (io.ReadCloser, error) {
	// TODO: Interesting options about logs are available, you can get them from user settings
	return d.cli.ContainerLogs(context.TODO(), containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Details:    false,
	})
}

func (d *docker) ContainerStart(containerID string) error {
	return d.cli.ContainerStart(context.TODO(), containerID, types.ContainerStartOptions{})
}

func (d *docker) ContainerStop(containerID string) error {
	return d.cli.ContainerStop(context.TODO(), containerID, nil)
}

func (d *docker) ContainerRemove(containerID string, removeForm *entities.ContainerRemoveForm) error {
	return d.cli.ContainerRemove(context.TODO(), containerID, types.ContainerRemoveOptions{
		RemoveVolumes: removeForm.RemoveVolumes,
		Force:         removeForm.Force,
	})
}

func (d *docker) ContainerRename(containerID, newName string) error {
	return d.cli.ContainerRename(context.TODO(), containerID, newName)
}
