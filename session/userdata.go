package session

import "github.com/arshamalh/dockeroller/entities"

type UserData struct {
	UserID           int64
	Scene            entities.Scene
	CurrentQuestion  entities.Question
	Containers       []*entities.Container
	Images           []*entities.Image
	CurrentContainer *entities.Container
	CurrentImage     *entities.Image
	Forms            *entities.Forms
	QuitChan         chan struct{}
}

func (d *UserData) SetScene(scene entities.Scene) {
	d.Scene = scene
	switch scene {
	case entities.SceneRenameContainer:
		d.CurrentQuestion = entities.QNewContainerName
	case entities.SceneRenameImage:
		d.CurrentQuestion = entities.QNewImageName
	}
}

func (d *UserData) GetScene() entities.Scene {
	return d.Scene
}

func (d *UserData) SetCurrentQuestion(question entities.Question) {
	d.CurrentQuestion = question
}

func (d *UserData) GetCurrentQuestion() entities.Question {
	return d.CurrentQuestion
}

func (d *UserData) SetCurrentContainer(container *entities.Container) {
	d.CurrentContainer = container
}

func (d *UserData) GetCurrentContainer() *entities.Container {
	return d.CurrentContainer
}

func (d *UserData) SetCurrentImage(image *entities.Image) {
	d.CurrentImage = image
}

func (d *UserData) GetCurrentImage() *entities.Image {
	return d.CurrentImage
}

func (d *UserData) GetContainers() []*entities.Container {
	return d.Containers
}

func (d *UserData) SetContainers(containers []*entities.Container) {
	d.Containers = containers
}

func (d *UserData) SetQuitChan(quitChan chan struct{}) {
	d.QuitChan = quitChan
}

func (d *UserData) GetQuitChan() chan<- struct{} {
	return d.QuitChan
}

func (d *UserData) GetImages() []*entities.Image {
	return d.Images
}

func (d *UserData) SetImages(images []*entities.Image) {
	d.Images = images
}

func (d *UserData) GetForms() *entities.Forms {
	return d.Forms
}

func (d *UserData) SetForms(forms *entities.Forms) {
	d.Forms = forms
}

func (d *UserData) SetContainerRemoveForm(force, removeVolumes bool) *entities.ContainerRemoveForm {
	if d.Forms == nil {
		d.Forms = &entities.Forms{
			ContainerRemove: &entities.ContainerRemoveForm{},
		}
	}
	d.Forms.ContainerRemove.Force = force
	d.Forms.ContainerRemove.RemoveVolumes = removeVolumes
	return d.Forms.ContainerRemove
}

func (d *UserData) GetContainerRemoveForm() *entities.ContainerRemoveForm {
	if d.Forms != nil && d.Forms.ContainerRemove != nil {
		return d.Forms.ContainerRemove
	}
	return nil
}

func (d *UserData) SetImageRemoveForm(force, pruneChildren bool) *entities.ImageRemoveForm {
	if d.Forms == nil {
		d.Forms = &entities.Forms{
			ImageRemove: &entities.ImageRemoveForm{},
		}
	}
	d.Forms.ImageRemove.Force = force
	d.Forms.ImageRemove.PruneChildren = pruneChildren
	return d.Forms.ImageRemove
}

func (d *UserData) GetImageRemoveForm() *entities.ImageRemoveForm {
	if d.Forms != nil && d.Forms.ImageRemove != nil {
		return d.Forms.ImageRemove
	}
	return nil
}
