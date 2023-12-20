package handlers

import (
	"context"
	"strconv"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/docker/docker/api/types/filters"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainerRename(ctx telebot.Context) error {
	if err := ctx.Respond(); err != nil {
		log.Gl.Error(err.Error())
	}
	userID := ctx.Chat().ID
	currentContainerIndex := ctx.Data()
	session := h.session.Get(userID)
	index, err := strconv.Atoi(currentContainerIndex)
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Send("wrong button clicked!")
	}
	containers := session.GetContainers()
	current := containers[index]
	session.SetScene(entities.SceneRenameContainer)
	session.SetCurrentContainer(current)

	return ctx.Edit(
		msgs.ContainerNewNameInput,
		keyboards.ContainerBack(index),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerRenameTextHandler(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	container := session.GetCurrentContainer()
	if container == nil {
		return ctx.Edit(
			"you're lost!, please /start again",
			keyboards.ContainerBack(0),
			telebot.ModeMarkdownV2,
		)
	}

	newName := ctx.Text()
	if err := h.docker.ContainerRename(container.ID, newName); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Edit(
			"we cannot rename this container",
			keyboards.ContainerBack(0),
			telebot.ModeMarkdownV2,
		)
	}

	containers := h.docker.ContainersList(context.TODO(), filters.Args{})
	session.SetContainers(containers)

	return ctx.Send(
		msgs.FmtContainerRenamed(container.Name, newName),
		keyboards.ContainerBack(0),
		telebot.ModeMarkdownV2,
	)
}
