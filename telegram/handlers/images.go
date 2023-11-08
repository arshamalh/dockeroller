package handlers

import (
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"gopkg.in/telebot.v3"
)

func (h *handler) ImagesHandler(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	images := h.docker.ImagesList()
	h.session.SetImages(userID, images)
	current := images[0]
	return ctx.Send(
		msgs.FmtImage(current),
		keyboards.ContainersList(0, false),
		telebot.ModeMarkdownV2,
	)
}
