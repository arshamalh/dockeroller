package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types"
)

func (d *docker) ImagesList() (images []*entities.Image) {
	rawImages, _ := d.cli.ImageList(context.TODO(), types.ImageListOptions{All: true})
	for _, rawImg := range rawImages {
		status := d.getImageStatus(context.TODO(), rawImg)
		images = append(images, &entities.Image{
			ID:        rawImg.ID,
			Size:      rawImg.Size,
			Tags:      rawImg.RepoTags,
			Status:    entities.ImageStatus(status),
			CreatedAt: fmt.Sprint(time.Unix(rawImg.Created, 0).Format("2006-01-02 15:04:05")),
		})
	}
	return
}

func (d *docker) ImageTag(ctx context.Context, imageID, newTag string) error {
	return d.cli.ImageTag(ctx, imageID, newTag)
}

func (d *docker) ImageRemove(ctx context.Context, imageID string, force, pruneChildren bool) error {
	_, err := d.cli.ImageRemove(ctx, imageID,
		types.ImageRemoveOptions{
			Force:         force,
			PruneChildren: pruneChildren,
		},
	)
	return err
}

func (d *docker) getImageStatus(ctx context.Context, image types.ImageSummary) (status string) {
	if len(image.RepoTags) == 0 {
		status = string(entities.ImageStatusUnUsedDangling)
		return
	}

	if image.RepoTags[0] == "<none>:<none>" {
		status = string(entities.ImageStatusUnUsedDangling)
		return
	}
	containers, _ := d.cli.ContainerList(ctx, types.ContainerListOptions{})
	if len(containers) == 0 {
		status = string(entities.ImageStatusUnUsed)
		return
	}
	newImgs := make(map[string][]string)
	for _, cont := range containers {
		if cont.ImageID != image.ID {
			status = string(entities.ImageStatusUnUsed)
		} else {
			newSlice := newImgs[image.ID]
			if newSlice == nil {
				newSlice = make([]string, 0)
			}
			newSlice = append(newSlice, cont.ID)
			newImgs[image.ID] = newSlice

			status = string(entities.ImageStatusInUse)
		}
	}

	return
}
