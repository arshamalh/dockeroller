package docker

import (
	"context"

	"github.com/arshamalh/dockeroller/models"
	"github.com/docker/docker/api/types"
)

func (d docker) ImagesList() (images []*models.Image) {
	raw_images, _ := d.cli.ImageList(context.TODO(), types.ImageListOptions{All: true})
	for _, rimg := range raw_images {
		images = append(images, &models.Image{
			ID:   rimg.ID,
			Size: rimg.Size,
			Tags: rimg.RepoTags,
		})
	}
	return
}
