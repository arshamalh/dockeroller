package handlers_test

import (
	"context"
	"testing"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/session"
	"github.com/arshamalh/dockeroller/telegram/handlers"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/docker/docker/api/types/filters"
	"github.com/go-faker/faker/v4"
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

	containersList := fakeContainers(t, 3)
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

func fakeContainers(t *testing.T, howMany int) []*entities.Container {
	containersList := make([]*entities.Container, howMany)
	for i, _ := range containersList {
		container := new(entities.Container)
		if err := faker.FakeData(container); err != nil {
			t.Error("can't make stubs", err)
		}
		containersList[i] = container
	}
	return containersList
}
