package handlers

import (
	"context"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"gopkg.in/telebot.v3"
)

// ContainerRename Handler is called where rename button is clicked
// It asks new name from user and waits for input
func (h *handler) ContainerRename(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	_, index := session.GetCurrentContainer()
	session.SetScene(entities.SceneRenameContainer)

	h.EmptyResponder(ctx)
	return ctx.Edit(
		msgs.ContainerNewNameInput,
		keyboards.ContainerBack(index),
	)
}

func (h *handler) ContainerRenameTextHandler(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	container, index := session.GetCurrentContainer()
	if container == nil {
		return ctx.Edit(
			"you're lost!, please /start again",
			keyboards.ContainerBack(index),
		)
	}

	newName := ctx.Text()
	if err := h.docker.ContainerRename(context.TODO(), container.ID, newName); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Edit(
			"we cannot rename this container",
			keyboards.ContainerBack(index),
		)
	}

	container.Name = newName
	session.SetCurrentContainer(container, index)

	return ctx.Send(
		msgs.FmtContainerRenamed(container.Name, newName),
		keyboards.ContainerBack(index),
	)
}
