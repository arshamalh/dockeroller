package entities

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
	UsedBy []*Container
}
