package models

type UserData struct {
	ContainerRemoveForm *ContainerRemoveForm
}

type ContainerRemoveForm struct {
	Force         bool
	RemoveVolumes bool
}
