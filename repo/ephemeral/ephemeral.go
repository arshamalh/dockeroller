package ephemeral

import (
	"github.com/arshamalh/dockeroller/models"
)

type userData struct {
	Containers []*models.Container
	Images     []*models.Image
	QuitChan   chan struct{}
}

type ephemeral struct {
	usersData map[int64]userData
}

func New() *ephemeral {
	return &ephemeral{
		usersData: make(map[int64]userData),
	}
}

func (e *ephemeral) GetContainers(userID int64) []*models.Container {
	return e.usersData[userID].Containers
}

func (e *ephemeral) SetContainers(userID int64, containers []*models.Container) {
	uData := e.usersData[userID]
	uData.Containers = containers
	e.usersData[userID] = uData
}

func (e *ephemeral) SetQuitChan(userID int64, quitChan chan struct{}) {
	uData := e.usersData[userID]
	uData.QuitChan = quitChan
	e.usersData[userID] = uData
}

func (e *ephemeral) GetQuitChan(userID int64) chan<- struct{} {
	return e.usersData[userID].QuitChan
}

func (e *ephemeral) GetImages(userID int64) []*models.Image {
	return e.usersData[userID].Images // TODO: Not safe!! what if userID was not valid?
}

func (e *ephemeral) SetImages(userID int64, images []*models.Image) {
	uData := e.usersData[userID]
	uData.Images = images
	e.usersData[userID] = uData
}

func (e *ephemeral) GetUserData() {

}

func (e *ephemeral) SetUserData() {

}
