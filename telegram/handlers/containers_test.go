package handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/session"
	"github.com/arshamalh/dockeroller/telegram/handlers"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/docker/docker/api/types/filters"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v3"
)

//go:generate mockgen -destination mock_telebot_context_test.go -package handlers_test gopkg.in/telebot.v3 Context
//go:generate mockgen -destination mock_docker_test.go -package handlers_test github.com/arshamalh/dockeroller/docker Docker
// TODO: restructure the session and make its mocks here

func TestContainersList(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	docker := NewMockDocker(ctrl)
	ctx := NewMockContext(ctrl)
	ssn := session.New()

	containersList := fakeContainers(3)
	// pp.Print(containersList)

	offBot, err := telebot.NewBot(telebot.Settings{Offline: true})
	assert.Nil(err, "can't make offline bot")

	handler := handlers.NewHandler(offBot, docker, ssn)

	t.Run("no containers", func(t *testing.T) {
		ctx.EXPECT().Respond(msgs.NoContainer).Return(nil)
		ctx.EXPECT().Chat().Return(&telebot.Chat{ID: 467})
		docker.EXPECT().ContainersList(context.TODO(), filters.Args{}).Return([]*entities.Container{})

		err := handler.ContainersList(ctx)

		assert.Nil(err)
	})

	t.Run("picking first container and formatting", func(t *testing.T) {
		docker.EXPECT().ContainersList(context.TODO(), filters.Args{}).Return(containersList)
		ctx.EXPECT().Chat().Return(&telebot.Chat{ID: 467})

		ctx.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(func(receivedContainer interface{}, opts ...interface{}) error {
			formattedContainer := msgs.FmtContainer(containersList[0])
			assert.Equal(formattedContainer, receivedContainer)
			return nil
		})

		err := handler.ContainersList(ctx)

		assert.Nil(err)
	})

	t.Run("filling the session", func(t *testing.T) {})
}

func TestContainersNav(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	docker := NewMockDocker(ctrl)
	ctx := NewMockContext(ctrl)
	ssn := session.New()

	const NumberOfContainers = 3
	containersList := fakeContainers(NumberOfContainers)
	// pp.Print(containersList)

	offBot, err := telebot.NewBot(telebot.Settings{Offline: true})
	assert.Nil(err, "can't make offline bot")

	handler := handlers.NewHandler(offBot, docker, ssn)

	t.Run("next container when there is a next", func(t *testing.T) {
		ctx.EXPECT().Respond().Return(nil)
		ctx.EXPECT().Chat().Return(&telebot.Chat{ID: 467})
		ctx.EXPECT().Data().Return("1")

		docker.EXPECT().ContainersList(context.TODO(), filters.Args{}).Return(containersList)

		editMockedFunc := func(receivedContainer interface{}, opts ...interface{}) error {
			formattedContainer := msgs.FmtContainer(containersList[1])
			assert.Equal(formattedContainer, receivedContainer)
			return nil
		}

		ctx.EXPECT().Edit(gomock.Any(), gomock.Any()).DoAndReturn(editMockedFunc)

		err := handler.ContainersNavBtn(ctx)

		assert.Nil(err)
	})

	t.Run("previous container of the first, look at last", func(t *testing.T) {
		ctx.EXPECT().Respond().Return(nil)
		ctx.EXPECT().Chat().Return(&telebot.Chat{ID: 467})
		ctx.EXPECT().Data().Return("-1")

		docker.EXPECT().ContainersList(context.TODO(), filters.Args{}).Return(containersList)

		editMockedFunc := func(receivedContainer interface{}, opts ...interface{}) error {
			formattedContainer := msgs.FmtContainer(containersList[NumberOfContainers-1])
			assert.Equal(formattedContainer, receivedContainer)
			return nil
		}

		ctx.EXPECT().Edit(gomock.Any(), gomock.Any()).DoAndReturn(editMockedFunc)

		err := handler.ContainersNavBtn(ctx)

		assert.Nil(err)
	})

	t.Run("filling the session", func(t *testing.T) {})
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
		ranInt := faker.IntBetween(0, len(allContainerStates)-1)
		containersList[i] = &entities.Container{
			ID:     faker.UUID().V4(),
			Status: faker.Lorem().Sentence(10),
			Name:   faker.Person().FirstName(),
			Image: fmt.Sprintf(
				"%s:%f-%s%f",
				faker.Address().City(), faker.Float32(2, 10, 99),
				faker.Blood().Name(), faker.Float32(2, 10, 99),
			),
			State: allContainerStates[ranInt],
		}
	}
	return containersList
}
