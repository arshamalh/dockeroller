package docker

import (
	"context"

	"github.com/arshamalh/dockeroller/models"
	"github.com/docker/docker/api/types"
)

func (d docker) ContainersList() (containers []*models.Container) {
	raw_containers, _ := d.cli.ContainerList(context.TODO(), types.ContainerListOptions{All: true})
	for _, rcont := range raw_containers {
		containers = append(containers, &models.Container{
			Name:   rcont.Names[0],
			Image:  rcont.Image,
			Status: rcont.Status,
		})
	}
	return
}
