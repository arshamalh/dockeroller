package models

type ImageStatus string

const (
	ImageStatusInUse          ImageStatus = "In use"
	ImageStatusUnUsed         ImageStatus = "Un used"
	ImageStatusUnUsedDangling ImageStatus = "Un used (dangling)"
)

type Image struct {
	ID        string
	Size      int64
	Tags      []string
	Status    ImageStatus
	CreatedAt string
}
