package entities

import (
	"fmt"

	"github.com/arshamalh/dockeroller/tools"
)

type ImageStatus string

const (
	ImageStatusInUse          ImageStatus = "In use"
	ImageStatusUnUsed         ImageStatus = "Unused"
	ImageStatusUnUsedDangling ImageStatus = "Dangling"
)

type Image struct {
	ID        string
	Size      int64
	Tags      []string
	Status    ImageStatus
	CreatedAt string
	// A list of Containers using this image
	UsedBy     []*Container
	RemoveForm *ImageRemoveForm
}

func (img Image) String() string {
	size := tools.SizeToHumanReadable(img.Size)
	return fmt.Sprintf("%s - %s - %s - created at: %s", img.ID, size, img.Status, img.CreatedAt)
}

func (img *Image) ShortID() string {
	return img.ID[:LEN_IMG_TRIM]
}
