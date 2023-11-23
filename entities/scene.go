package entities

type Scene int

const (
	SceneRenameContainer Scene = iota + 1
	SceneRenameImage
)

type Question int

const (
	QNewContainerName Question = iota + 1
)

const (
	QNewImageName Question = iota + 1
)

func (q Question) NextQuestion() Question {
	return q + 1
}
