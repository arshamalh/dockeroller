package docker

import (
	"context"
	"io"

	"github.com/arshamalh/dockeroller/models"
	"github.com/docker/docker/api/types"
)

func (d docker) ContainersList() (containers []*models.Container) {
	raw_containers, _ := d.cli.ContainerList(context.TODO(), types.ContainerListOptions{All: true})
	for _, rcont := range raw_containers {
		containers = append(containers, &models.Container{
			ID:     rcont.ID,
			Name:   rcont.Names[0],
			Image:  rcont.Image,
			Status: rcont.Status,
		})
	}
	return
}

func (d docker) ContainerStats(containerID string) (io.ReadCloser, error) {
	stats, err := d.cli.ContainerStats(context.TODO(), containerID, true)
	return stats.Body, err
}
