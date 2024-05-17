package handlers

import (
	"context"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"gopkg.in/telebot.v3"
)

func (h *handler) ImageTag(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index := tools.Str2Int(ctx.Data())
	h.session.Get(userID).SetScene(entities.SceneRenameImage)

	h.EmptyResponder(ctx)
	return ctx.Edit(
		msgs.ImageNewNameInput,
		keyboards.ImageBack(index),
	)
}

func (h *handler) ImageTagTextHandler(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	image, index := session.GetCurrentImage()
	if image == nil {
		return ctx.Send(
			"you're lost!, please /start again",
			keyboards.ImageBack(0),
		)
	}

	newTag := ctx.Text()
	if err := h.docker.ImageTag(context.TODO(), image.ID, newTag); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Send(
			msgs.FmtMono(msgs.InvalidImageTag),
		)
	}

	return ctx.Send(
		msgs.FmtImageTagged(image.ID, newTag),
		keyboards.ImageBack(index),
	)
}
