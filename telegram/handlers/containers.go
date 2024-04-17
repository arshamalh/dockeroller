package handlers

import (
	"context"
	"strconv"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"github.com/docker/docker/api/types/filters"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainersList(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	containers := h.docker.ContainersList(context.TODO(), filters.Args{})
	session.SetContainers(containers)
	if len(containers) == 0 {
		return ctx.Respond(msgs.NoContainer)
	}
	current := containers[0]
	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainersNavBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}

	session := h.session.Get(userID)
	containers := h.docker.ContainersList(context.TODO(), filters.Args{})
	session.SetContainers(containers)
	if len(containers) == 0 {
		return ctx.Respond(msgs.NoContainer)
	}
	index = tools.Indexer(index, len(containers))
	current := containers[index]

	err = ctx.Respond()
	if err != nil {
		log.Gl.Error(err.Error())
	}
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
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
	current := session.GetContainer(index)

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, current.IsOn()),
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
	current := session.GetContainer(index)
	if err := h.docker.ContainerStart(current.ID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.CannotStartTheContainer)
	}

	current, err = h.docker.GetContainer(current.ID)
	if err != nil {
		return ctx.Respond(msgs.StartedButUnavailableCurrentState)
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
	current := session.GetContainer(index)
	if err := h.docker.ContainerStop(current.ID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.CannotStopTheContainer)
	}

	current, err = h.docker.GetContainer(current.ID)
	if err != nil {
		return ctx.Respond(msgs.StoppedButUnavailableCurrentState)
	}

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, false),
		telebot.ModeMarkdownV2,
	)
}
