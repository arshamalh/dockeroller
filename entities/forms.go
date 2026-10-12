package entities

// ContainerRemoveForm holds the options of the container remove form.
type ContainerRemoveForm struct {
	Force         bool
	RemoveVolumes bool
}

// ImageRemoveForm holds the options of the image remove form.
type ImageRemoveForm struct {
	Force         bool
	PruneChildren bool
}
