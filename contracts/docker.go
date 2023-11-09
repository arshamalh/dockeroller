package contracts

import (
	"io"

	"github.com/arshamalh/dockeroller/models"
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
