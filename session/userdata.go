package session

import "github.com/arshamalh/dockeroller/entities"

type UserData struct {
	userID                int64
	scene                 entities.Scene
	currentContainer      *entities.Container
	currentContainerIndex int
	currentImage          *entities.Image
	currentImageIndex     int
	quitChan              chan struct{}
}

func (d *UserData) ID() int64 {
	return d.userID
}

func (d *UserData) SetScene(scene entities.Scene) {
	d.scene = scene
}

func (d *UserData) GetScene() entities.Scene {
	return d.scene
}

func (d *UserData) SetCurrentContainer(container *entities.Container, index int) {
	if container.RemoveForm == nil {
		container.RemoveForm = &entities.ContainerRemoveForm{}
	}
	d.currentContainer = container
	d.currentContainerIndex = index
}

func (d *UserData) GetCurrentContainer() (*entities.Container, int) {
	return d.currentContainer, d.currentContainerIndex
}

func (d *UserData) SetCurrentImage(image *entities.Image, index int) {
	if image.RemoveForm == nil {
		image.RemoveForm = &entities.ImageRemoveForm{}
	}
	d.currentImage = image
	d.currentImageIndex = index
}

func (d *UserData) GetCurrentImage() (*entities.Image, int) {
	return d.currentImage, d.currentImageIndex
}

func (d *UserData) SetQuitChan(quitChan chan struct{}) {
	d.quitChan = quitChan
}

func (d *UserData) GetQuitChan() chan<- struct{} {
	return d.quitChan
}
