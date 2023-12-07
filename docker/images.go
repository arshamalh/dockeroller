package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (d *docker) ImagesList() (images []*entities.Image) {
	rawImages, _ := d.cli.ImageList(context.TODO(), types.ImageListOptions{All: true})
	for _, rawImg := range rawImages {
		status, containers := d.getImageStatus(context.TODO(), rawImg)
		images = append(images, &entities.Image{
			ID:        rawImg.ID,
			Size:      rawImg.Size,
			Tags:      rawImg.RepoTags,
			Status:    entities.ImageStatus(status),
			UsedBy:    containers,
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

// Returns whether an image is dangling or used by containers,
// returns a list of ContainerIDs as the second argument, or nil if the Image is dangling or unused
func (d *docker) getImageStatus(ctx context.Context, image types.ImageSummary) (entities.ImageStatus, []*entities.Container) {
	if len(image.RepoTags) == 0 {
		return entities.ImageStatusUnUsedDangling, nil
	}

	// Compatibility with older docker daemons
	if image.RepoTags[0] == "<none>:<none>" {
		return entities.ImageStatusUnUsedDangling, nil
	}

	containers := d.ContainersList(ctx, filters.NewArgs(filters.Arg("ancestor", image.ID)))
	if len(containers) == 0 {
		return entities.ImageStatusUnUsed, nil
	}

	return entities.ImageStatusInUse, containers
}
