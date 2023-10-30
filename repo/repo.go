package repo

import "github.com/arshamalh/dockeroller/models"

type Session interface {
	GetContainers(userID int64) []*models.Container
	SetContainers(userID int64, containers []*models.Container)
	SetQuitChan(userID int64, quitChan chan struct{})
	GetQuitChan(userID int64) chan<- struct{}
	GetImages(userID int64) []*models.Image
	SetImages(userID int64, images []*models.Image)
	GetUserData()
	SetUserData()
}
