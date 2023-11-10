package docker

import (
	"context"
	"fmt"
	"github.com/arshamalh/dockeroller/models"
	"github.com/docker/docker/api/types"
	"time"
)

func (d docker) ImagesList() (images []*models.Image) {
	raw_images, _ := d.cli.ImageList(context.TODO(), types.ImageListOptions{All: true})
	for _, rimg := range raw_images {
		status := d.getImageStatus(context.TODO(), rimg)
		images = append(images, &models.Image{
			ID:        rimg.ID,
			Size:      rimg.Size,
			Tags:      rimg.RepoTags,
			Status:    models.ImageStatus(status),
			CreatedAt: fmt.Sprint(time.Unix(rimg.Created, 0).Format("2006-01-02 15:04:05")),
		})
	}
	return
}

func (d docker) getImageStatus(ctx context.Context, image types.ImageSummary) (status string) {
	if len(image.RepoTags) == 0 {
		status = string(models.ImageStatusUnUsedDangling)
	}

	containers, _ := d.cli.ContainerList(ctx, types.ContainerListOptions{})
	newImgs := make(map[string][]string)
	for _, cont := range containers {
		if cont.ImageID != image.ID {
			status = string(models.ImageStatusUnUsed)
		} else {
			newSlice := newImgs[image.ID]
			if newSlice == nil {
				newSlice = make([]string, 0)
			}
			newSlice = append(newSlice, cont.ID)
			newImgs[image.ID] = newSlice

			status = string(models.ImageStatusInUse)
		}
	}

	return
}
