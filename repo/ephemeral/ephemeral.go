package ephemeral

import (
	"github.com/arshamalh/dockeroller/models"
)

type data struct {
	Scene            models.Scene
	CurrentQuestion  models.Question
	Containers       []*models.Container
	Images           []*models.Image
	CurrentContainer *models.Container
	CurrentImage     *models.Image
	UserData         *models.UserData
	QuitChan         chan struct{}
}

type ephemeral struct {
	data map[int64]data
}

func New() *ephemeral {
	return &ephemeral{
		data: make(map[int64]data),
	}
}

func (e *ephemeral) SetScene(userID int64, scene models.Scene) {
	uData := e.data[userID]
	uData.Scene = scene
	switch scene {
	case models.SceneRenameContainer:
		uData.CurrentQuestion = models.QNewContainerName
	case models.SceneRenameImage:
		uData.CurrentQuestion = models.QNewImageName
	}
	e.data[userID] = uData
}

func (e *ephemeral) GetScene(userID int64) models.Scene {
	return e.data[userID].Scene
}

func (e *ephemeral) SetCurrentQuestion(userID int64, question models.Question) {
	uData := e.data[userID]
	uData.CurrentQuestion = question
	e.data[userID] = uData
}

func (e *ephemeral) GetCurrentQuestion(userID int64) models.Question {
	return e.data[userID].CurrentQuestion
}

func (e *ephemeral) SetCurrentContainer(userID int64, container *models.Container) {
	uData := e.data[userID]
	uData.CurrentContainer = container
	e.data[userID] = uData
}

func (e *ephemeral) GetCurrentContainer(userID int64) *models.Container {
	return e.data[userID].CurrentContainer
}

func (e *ephemeral) SetCurrentImage(userID int64, image *models.Image) {
	uData := e.data[userID]
	uData.CurrentImage = image
	e.data[userID] = uData
}

func (e *ephemeral) GetCurrentImage(userID int64) *models.Image {
	return e.data[userID].CurrentImage
}

func (e *ephemeral) GetContainers(userID int64) []*models.Container {
	return e.data[userID].Containers
}

func (e *ephemeral) SetContainers(userID int64, containers []*models.Container) {
	uData := e.data[userID]
	uData.Containers = containers
	e.data[userID] = uData
}

func (e *ephemeral) SetQuitChan(userID int64, quitChan chan struct{}) {
	uData := e.data[userID]
	uData.QuitChan = quitChan
	e.data[userID] = uData
}

func (e *ephemeral) GetQuitChan(userID int64) chan<- struct{} {
	return e.data[userID].QuitChan
}

func (e *ephemeral) GetImages(userID int64) []*models.Image {
	return e.data[userID].Images // TODO: Not safe!! what if userID was not valid?
}

func (e *ephemeral) SetImages(userID int64, images []*models.Image) {
	uData := e.data[userID]
	uData.Images = images
	e.data[userID] = uData
}

func (e *ephemeral) GetUserData(userID int64) *models.UserData {
	return e.data[userID].UserData
}

func (e *ephemeral) SetUserData(userID int64, userData *models.UserData) {
	uData := e.data[userID]
	uData.UserData = userData
	e.data[userID] = uData
}

func (e *ephemeral) SetContainerRemoveForm(userID int64, force, removeVolumes bool) *models.ContainerRemoveForm {
	uData := e.data[userID]
	if uData.UserData == nil {
		uData.UserData = &models.UserData{
			ContainerRemoveForm: &models.ContainerRemoveForm{},
		}
	}
	uData.UserData.ContainerRemoveForm.Force = force
	uData.UserData.ContainerRemoveForm.RemoveVolumes = removeVolumes
	e.data[userID] = uData
	return uData.UserData.ContainerRemoveForm
}

func (e *ephemeral) GetContainerRemoveForm(userID int64) *models.ContainerRemoveForm {
	if e.data[userID].UserData != nil && e.data[userID].UserData.ContainerRemoveForm != nil {
		return e.data[userID].UserData.ContainerRemoveForm
	}
	return nil
}

func (e *ephemeral) SetImageRemoveForm(userID int64, force, pruneChildren bool) *models.ImageRemoveForm {
	uData := e.data[userID]
	if uData.UserData == nil {
		uData.UserData = &models.UserData{
			ImageRemoveForm: &models.ImageRemoveForm{},
		}
	}
	uData.UserData.ImageRemoveForm.Force = force
	uData.UserData.ImageRemoveForm.PruneChildren = pruneChildren
	e.data[userID] = uData
	return uData.UserData.ImageRemoveForm
}

func (e *ephemeral) GetImageRemoveForm(userID int64) *models.ImageRemoveForm {
	if e.data[userID].UserData != nil && e.data[userID].UserData.ImageRemoveForm != nil {
		return e.data[userID].UserData.ImageRemoveForm
	}
	return nil
}
