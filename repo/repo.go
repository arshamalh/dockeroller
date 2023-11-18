package repo

import "github.com/arshamalh/dockeroller/models"

type Session interface {
	SetScene(userID int64, scene models.Scene)
	GetScene(userID int64) models.Scene

	SetCurrentQuestion(userID int64, question models.Question)
	GetCurrentQuestion(userID int64) models.Question

	SetCurrentContainer(userID int64, container *models.Container)
	GetCurrentContainer(userID int64) *models.Container

	SetCurrentImage(userID int64, image *models.Image)
	GetCurrentImage(userID int64) *models.Image

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

	SetImageRemoveForm(userID int64, force, removeVolumes bool) *models.ImageRemoveForm
	GetImageRemoveForm(userID int64) *models.ImageRemoveForm
}
