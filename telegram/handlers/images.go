package handlers

import (
	"strconv"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"gopkg.in/telebot.v3"
)

func (h *handler) ImagesList(ctx telebot.Context) error {
	ctx.Respond()
	userID := ctx.Chat().ID
	images := h.updateImagesList(userID)
	current := images[0]
	return ctx.Send(
		msgs.FmtImage(current),
		keyboards.ImagesList(0),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ImagesNavBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}

	images := h.updateImagesList(userID)
	if len(images) == 0 {
		return ctx.Respond(&telebot.CallbackResponse{Text: "There is either no images or you should run /images again!"})
	}
	index = tools.Indexer(index, len(images))
	current := images[index]
	err = ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImagesList(index),
		telebot.ModeMarkdownV2,
	)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	return ctx.Respond()
}

func (h *handler) ImagesBackBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	if quitChan := session.GetQuitChan(); quitChan != nil {
		quitChan <- struct{}{}
	}
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetImages()[index]
	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImagesList(index),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) updateImagesList(userID int64) []*entities.Image {
	images := h.docker.ImagesList()
	session := h.session.Get(userID)
	session.SetImages(images)
	return images
}
