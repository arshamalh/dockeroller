package repo

import "github.com/arshamalh/dockeroller/models"

type Session interface {
	SetContainers(userID int64, containers []*models.Container)
	GetContainers(userID int64) []*models.Container

	SetQuitChan(userID int64, quitChan chan struct{})
	GetQuitChan(userID int64) chan<- struct{}

	SetImages(userID int64, images []*models.Image)
	GetImages(userID int64) []*models.Image

	SetUserData(userID int64, userData *models.UserData)
	GetUserData(userID int64) *models.UserData

	SetContainerRemoveForm(userID int64, force, removeVolumes bool) *models.ContainerRemoveForm
	GetContainerRemoveForm(userID int64) *models.ContainerRemoveForm
}
