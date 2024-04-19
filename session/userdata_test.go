package session_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/session"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestUserData(t *testing.T) {
	assert := assert.New(t)

	t.Run("set and get current container", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)
		fakeContainers := fakeContainers(1)

		userData.SetCurrentContainer(fakeContainers[0], 0)
		userCurrentContainer, index := userData.GetCurrentContainer()

		assert.NotNil(userCurrentContainer)
		assert.Equal(0, index)
	})

	t.Run("set and get current image", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)
		fakeImages := fakeImages(3)

		userData.SetCurrentImage(fakeImages[1], 1)
		userCurrentImage, _ := userData.GetCurrentImage()

		assert.NotNil(userCurrentImage)
	})

	t.Run("set and get scene", func(t *testing.T) {
		s := session.New()
		userData := s.Get(2)

		userData.SetScene(entities.SceneRenameImage)
		scene := userData.GetScene()
		assert.NotNil(scene)

		userData.SetScene(entities.SceneRenameContainer)
		scene = userData.GetScene()
		assert.NotNil(scene)
	})
}

func fakeContainers(howMany int) []*entities.Container {
	faker := faker.New()
	allContainerStates := []entities.ContainerState{
		entities.ContainerStateCreated,
		entities.ContainerStateRunning,
		entities.ContainerStateDead,
		entities.ContainerStateExited,
		entities.ContainerStatePaused,
	}

	containersList := make([]*entities.Container, howMany)

	for i := range containersList {
		randInt := faker.IntBetween(0, len(allContainerStates)-1)
		containersList[i] = &entities.Container{
			ID:     faker.UUID().V4(),
			Status: faker.Lorem().Sentence(10),
			Name:   faker.Person().FirstName(),
			Image: fmt.Sprintf(
				"%s:%f-%s%f",
				faker.Address().City(), faker.Float32(2, 10, 99),
				faker.Blood().Name(), faker.Float32(2, 10, 99),
			),
			State: allContainerStates[randInt],
		}
	}
	return containersList
}

func fakeImages(howMany int) []*entities.Image {
	faker := faker.New()

	imagesList := make([]*entities.Image, howMany)

	for i := range imagesList {
		imagesList[i] = &entities.Image{
			ID: faker.UUID().V4(),
			// TODO: image status is enum
			Status:    entities.ImageStatus(faker.Lorem().Sentence(10)),
			Size:      faker.Int64(),
			CreatedAt: faker.Time().RFC1123(time.Now()),
			UsedBy:    fakeContainers(3),
		}
	}
	return imagesList
}
