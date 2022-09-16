package contracts

import (
	"io"

	"github.com/arshamalh/dockeroller/models"
)

type Docker interface {
	ContainersList() []*models.Container
	ImagesList() []*models.Image
	// ContainerLogs()
	ContainerStats(containerID string) (io.ReadCloser, error)
}
