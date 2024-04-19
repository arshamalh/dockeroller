package session

import "github.com/arshamalh/dockeroller/entities"

type UserData struct {
	userID                int64
	scene                 entities.Scene
	currentContainer      *entities.Container
	currentContainerIndex int
	currentImage          *entities.Image
	currentImageIndex     int
	forms                 *entities.Forms
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

func (d *UserData) SetImageRemoveForm(force, pruneChildren bool) *entities.ImageRemoveForm {
	if d.forms == nil {
		d.forms = &entities.Forms{
			ImageRemove: &entities.ImageRemoveForm{},
		}
	}
	d.forms.ImageRemove.Force = force
	d.forms.ImageRemove.PruneChildren = pruneChildren
	return d.forms.ImageRemove
}

func (d *UserData) GetImageRemoveForm() *entities.ImageRemoveForm {
	if d.forms != nil && d.forms.ImageRemove != nil {
		return d.forms.ImageRemove
	}
	return &entities.ImageRemoveForm{}
}
