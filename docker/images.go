package docker

import (
	"context"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

func (d *docker) ImagesList(ctx context.Context) ([]*entities.Image, error) {
	rawImages, err := d.cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	images := make([]*entities.Image, 0, len(rawImages))
	for _, rawImg := range rawImages {
		status, containers := d.getImageStatus(ctx, rawImg)
		images = append(images, &entities.Image{
			ID:        rawImg.ID,
			Size:      rawImg.Size,
			Tags:      rawImg.RepoTags,
			Status:    status,
			UsedBy:    containers,
			CreatedAt: time.Unix(rawImg.Created, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return images, nil
}

func (d *docker) ImageTag(ctx context.Context, imageID, newTag string) error {
	return d.cli.ImageTag(ctx, imageID, newTag)
}

func (d *docker) ImageRemove(ctx context.Context, imageID string, force, pruneChildren bool) error {
	_, err := d.cli.ImageRemove(ctx, imageID, image.RemoveOptions{
		Force:         force,
		PruneChildren: pruneChildren,
	})
	return err
}

// Returns whether an image is dangling or used by containers,
// returns a list of ContainerIDs as the second argument, or nil if the Image is dangling or unused
func (d *docker) getImageStatus(ctx context.Context, img image.Summary) (entities.ImageStatus, []*entities.Container) {
	if len(img.RepoTags) == 0 {
		return entities.ImageStatusUnUsedDangling, nil
	}
	// Compatibility with older docker daemons
	if img.RepoTags[0] == "<none>:<none>" {
		return entities.ImageStatusUnUsedDangling, nil
	}
	containers, err := d.ContainersList(ctx, filters.NewArgs(filters.Arg("ancestor", img.ID)))
	if err != nil || len(containers) == 0 {
		return entities.ImageStatusUnUsed, nil
	}
	return entities.ImageStatusInUse, containers
}
