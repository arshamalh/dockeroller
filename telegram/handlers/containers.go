package handlers

import (
	"context"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"github.com/docker/docker/api/types/filters"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainersList(ctx telebot.Context) error {
	containers, err := h.docker.ContainersList(context.TODO(), filters.Args{})
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.UnableToFetchContainers)
	}
	if len(containers) == 0 {
		return ctx.Respond(msgs.NoContainer)
	}
	current := containers[0]
	h.session.Get(ctx.Chat().ID).SetCurrentContainer(current, 0)

	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(current.ID, 0, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainersNavBtn(ctx telebot.Context) error {
	containers, err := h.docker.ContainersList(context.TODO(), filters.Args{})
	if err != nil {
		return ctx.Respond(msgs.UnableToFetchContainers)
	}

	if len(containers) == 0 {
		return ctx.Respond(msgs.NoContainer)
	}

	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	index := tools.Str2Int(ctx.Data())
	index = tools.Indexer(index, len(containers))
	current := containers[index]
	session.SetCurrentContainer(current, index)

	h.EmptyResponder(ctx)
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(current.ID, index, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainersBackBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	if quitChan := session.GetQuitChan(); quitChan != nil {
		quitChan <- struct{}{}
	}

	current, index := session.GetCurrentContainer()
	if current == nil {
		containers, err := h.docker.ContainersList(context.TODO(), filters.Args{})
		if err != nil {
			log.Gl.Error(err.Error())
			ctx.Respond(msgs.UnableToFetchContainers)
			return h.StartHandler(ctx)
		}
		if len(containers) == 0 {
			ctx.Respond(msgs.NoContainer)
			return h.StartHandler(ctx)
		}
		current = containers[0]
		session.SetCurrentContainer(current, index)
	}

	h.EmptyResponder(ctx)
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(current.ID, index, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerStart(ctx telebot.Context) error {
	containerID := ctx.Data()
	if err := h.docker.ContainerStart(context.TODO(), containerID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.CannotStartTheContainer)
	}

	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	current, index := session.GetCurrentContainer()

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(current.ID, index, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerStop(ctx telebot.Context) error {
	containerID := ctx.Data()
	if err := h.docker.ContainerStop(context.TODO(), containerID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.CannotStopTheContainer)
	}

	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	current, index := session.GetCurrentContainer()

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(current.ID, index, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}
