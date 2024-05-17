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

func (h *handler) ImagesList(ctx telebot.Context) error {
	images, err := h.docker.ImagesList(context.TODO(), filters.Args{})
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.UnableToFetchImages)
	}
	if len(images) == 0 {
		return ctx.Respond(msgs.NoImages)
	}

	current := images[0]
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	session.SetCurrentImage(current, 0)

	h.EmptyResponder(ctx)
	return ctx.Send(
		msgs.FmtImage(current),
		keyboards.ImagesList(0),
	)
}

func (h *handler) ImagesNavBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)

	images, err := h.docker.ImagesList(context.TODO(), filters.Args{})
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.UnableToFetchImages)
	}
	if len(images) == 0 {
		return ctx.Respond(msgs.NoImages)
	}

	index := tools.Str2Int(ctx.Data())
	index = tools.Indexer(index, len(images))
	current := images[index]
	session.SetCurrentImage(current, index)

	h.EmptyResponder(ctx)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImagesList(index),
	)
}

func (h *handler) ImagesBackBtn(ctx telebot.Context) error {
	images, err := h.docker.ImagesList(context.TODO(), filters.Args{})
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.UnableToFetchImages)
	}
	if len(images) == 0 {
		return ctx.Respond(msgs.NoImages)
	}

	index := tools.Str2Int(ctx.Data())
	index = tools.Indexer(index, len(images))
	current := images[index]
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	session.SetCurrentImage(current, index)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImagesList(index),
	)
}
