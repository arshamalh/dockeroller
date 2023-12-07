package handlers

import (
	"context"
	"strconv"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"github.com/docker/docker/api/types/filters"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainersList(ctx telebot.Context) error {
	ctx.Respond()
	userID := ctx.Chat().ID
	containers := h.updateContainersList(userID)
	if len(containers) == 0 {
		return ctx.Send("there is no container")
	}
	current := containers[0]
	containerIsOn := current.State == entities.ContainerStateRunning
	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, containerIsOn),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainersNavBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}

	containers := h.updateContainersList(userID)
	if len(containers) == 0 {
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "There is either no containers or you should run /containers again!",
			},
		)
	}
	index = tools.Indexer(index, len(containers))
	current := containers[index]

	containerIsOn := current.State == entities.ContainerStateRunning
	err = ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, containerIsOn),
		telebot.ModeMarkdownV2,
	)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	return ctx.Respond()
}

func (h *handler) ContainersBackBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	if quitChan := session.GetQuitChan(); quitChan != nil {
		quitChan <- struct{}{}
	}
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetContainers()[index]

	containerIsOn := current.State == entities.ContainerStateRunning
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, containerIsOn),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerStart(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetContainers()[index]
	if err := h.docker.ContainerStart(current.ID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "We cannot start the container!",
			},
		)
	}

	current, err = h.docker.GetContainer(current.ID)
	if err != nil {
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "Container started, but we're not able to show current state.",
			},
		)
	}
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, true),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerStop(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetContainers()[index]
	if err := h.docker.ContainerStop(current.ID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "We cannot stop the container!",
			},
		)
	}

	current, err = h.docker.GetContainer(current.ID)
	if err != nil {
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "Container stopped, but we're not able to show current state.",
			},
		)
	}

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, false),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) updateContainersList(userID int64) []*entities.Container {
	containers := h.docker.ContainersList(context.TODO(), filters.Args{})
	session := h.session.Get(userID)
	session.SetContainers(containers)
	return containers
}
