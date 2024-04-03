package handlers_test

import (
	"testing"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/telegram/handlers"
	"github.com/docker/docker/api/types/filters"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination mock_telebot_context_test.go -package handlers_test gopkg.in/telebot.v3 Context
//go:generate mockgen -destination mock_docker_test.go -package handlers_test github.com/arshamalh/dockeroller/docker Docker
// TODO: restructure the session and make its mocks here

func TestContainersList(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)

	docker := NewMockDocker(ctrl)
	ctx := NewMockContext(ctrl)

	handler := handlers.NewHandler(nil, docker, nil)

	t.Run("no containers", func(t *testing.T) {
		// Arrange: setting what should be returned from the mocks
		// ctx.EXPECT().Respond().Do()
		docker.EXPECT().ContainersList(ctx, filters.Args{}).DoAndReturn([]*entities.Container{})

		// Act
		err := handler.ContainersList(ctx)

		// Assert
		assert.Nil(err)
	})

	t.Run("picking first container and formatting", func(t *testing.T) {

	})

	t.Run("filling the session", func(t *testing.T) {})
}
