package msgs

import (
	"strings"

	"github.com/arshamalh/dockeroller/models"
)

func FormatContainer(container *models.Container) string {
	response := strings.NewReplacer(
		"{Name}", container.Name,
		"{Image}", container.Image,
		"{Status}", container.Status,
	).Replace(Container)
	response = FmtMono(response)
	return response
}
