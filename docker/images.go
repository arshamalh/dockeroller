package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types"
)

func (d *docker) ImagesList() (images []*entities.Image) {
	raw_images, _ := d.cli.ImageList(context.TODO(), types.ImageListOptions{All: true})
	for _, raw_img := range raw_images {
		status := d.getImageStatus(context.TODO(), raw_img)
		images = append(images, &entities.Image{
			ID:        raw_img.ID,
			Size:      raw_img.Size,
			Tags:      raw_img.RepoTags,
			Status:    entities.ImageStatus(status),
			CreatedAt: fmt.Sprint(time.Unix(raw_img.Created, 0).Format("2006-01-02 15:04:05")),
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
	if image.RepoTags[0] == "<none>:<none>" {
		status = string(entities.ImageStatusUnUsedDangling)
		return
	}

	if image.Containers == 0 {
		status = string(entities.ImageStatusUnUsed)
		return
	}

	containers, _ := d.cli.ContainerList(ctx, types.ContainerListOptions{All: true, Latest: true})
	for _, cont := range containers {
		if cont.ImageID == image.ID {
			status = string(entities.ImageStatusInUse)
			return
		}
	}

	return string(entities.ImageStatusUnUsed)
}
