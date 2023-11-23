package handlers

import (
	"context"
	"strconv"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"gopkg.in/telebot.v3"
)

func (h *handler) ImageTag(ctx telebot.Context) error {
	if err := ctx.Respond(); err != nil {
		log.Gl.Error(err.Error())
	}
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	currentImageIndex := ctx.Data()
	index, err := strconv.Atoi(currentImageIndex)
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Send("wrong button clicked!")
	}
	images := session.GetImages()
	current := images[index]
	session.SetScene(entities.SceneRenameImage)
	session.SetCurrentImage(current)

	return ctx.Edit(
		msgs.ImageNewNameInput,
		keyboards.ImageBack(index),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ImageTagTextHandler(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	image := session.GetCurrentImage()
	if image == nil {
		return ctx.Edit(
			"you're lost!, please /start again",
			keyboards.ImageBack(0),
			telebot.ModeMarkdownV2,
		)
	}

	newTag := ctx.Text()
	if err := h.docker.ImageTag(context.TODO(), image.ID, newTag); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Edit(
			"we cannot rename this image",
			keyboards.ImageBack(0),
			telebot.ModeMarkdownV2,
		)
	}

	h.updateImagesList(userID)

	return ctx.Send(
		msgs.FmtImageTagged(image.ID, newTag),
		keyboards.ImageBack(0),
		telebot.ModeMarkdownV2,
	)
}
