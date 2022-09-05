package contracts

import "github.com/arshamalh/dockeroller/models"

type Docker interface {
	ContainersList() []*models.Container
	ImagesList() []*models.Image
}
