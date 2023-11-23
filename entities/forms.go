package entities

type Forms struct {
	ContainerRemove *ContainerRemoveForm
	ImageRemove     *ImageRemoveForm
}

type ContainerRemoveForm struct {
	Force         bool
	RemoveVolumes bool
}

type ImageRemoveForm struct {
	Force         bool
	PruneChildren bool
}
