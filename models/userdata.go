package models

type UserData struct {
	ContainerRemoveForm *ContainerRemoveForm
	ImageRemoveForm     *ImageRemoveForm
}

type ContainerRemoveForm struct {
	Force         bool
	RemoveVolumes bool
}

type ImageRemoveForm struct {
	Force         bool
	PruneChildren bool
}
