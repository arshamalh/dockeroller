package session_test

import (
	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/session"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserData(t *testing.T) {
	assert := assert.New(t)
	t.Run("set and get containers", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)
		fakeContainers := newFakeContainers()

		userData.SetContainers(fakeContainers)
		userContainers := userData.GetContainers()

		assert.NotNil(userContainers)
		for _, userContainer := range userContainers {
			assert.NotNil(userContainer)
		}
	})

	t.Run("set and get current container", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)
		fakeContainers := newFakeContainers()

		userData.SetCurrentContainer(fakeContainers[1])
		userCurrentContainer := userData.GetCurrentContainer()

		assert.NotNil(userCurrentContainer)
	})

	t.Run("set and get images", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)
		fakeImages := newFakeImages()

		userData.SetImages(fakeImages)
		userImages := userData.GetImages()

		assert.NotNil(userImages)
		for _, userImage := range userImages {
			assert.NotNil(userImage)
		}
	})

	t.Run("set and get current image", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)
		fakeImages := newFakeImages()

		userData.SetCurrentImage(fakeImages[1])
		userCurrentImage := userData.GetCurrentImage()

		assert.NotNil(userCurrentImage)
	})

	t.Run("set and get scene", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)

		userData.SetScene(entities.SceneRenameImage)
		scene := userData.GetScene()

		assert.NotNil(scene)
		assert.Equal(userData.CurrentQuestion, entities.QNewImageName)

		userData.SetScene(entities.SceneRenameContainer)
		scene = userData.GetScene()

		assert.NotNil(scene)
		assert.Equal(userData.CurrentQuestion, entities.QNewContainerName)
	})
}

func newFakeContainers() []*entities.Container {
	var containers []*entities.Container
	for i := 0; i < 3; i++ {
		containers = append(containers, &entities.Container{
			ID:     faker.UUIDDigit(),
			Name:   faker.Word(),
			State:  entities.ContainerStateRunning,
			Image:  faker.Word(),
			Status: faker.Word(),
		})
	}

	return containers
}

func newFakeImages() []*entities.Image {
	var images []*entities.Image
	for i := 0; i < 3; i++ {
		images = append(images, &entities.Image{
			ID:        faker.UUIDDigit(),
			Size:      12,
			Tags:      []string{faker.Word(), faker.WORD},
			Status:    entities.ImageStatusInUse,
			CreatedAt: faker.TIMESTAMP,
		})
	}

	return images
}
