package handlers

import (
	"context"
	"strconv"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"gopkg.in/telebot.v3"
)

// This handler doesn't remove the image,
// it's just a form for setting the options of removing a image in the later "Done" step.
func (h *handler) ImageRemoveForm(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.InvalidButton)
	}
	ctx.Respond(msgs.FillTheFormAndPressDone)

	current := session.GetImages()[index]
	imgRmForm := session.GetImageRemoveForm()
	if imgRmForm == nil {
		imgRmForm = session.SetImageRemoveForm(false, false)
	}

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(index, imgRmForm.Force, imgRmForm.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}

// Removes the image with specified options
func (h *handler) ImageRemoveDone(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetImages()[index]
	imgRmForm := session.GetImageRemoveForm()

	if err := h.docker.ImageRemove(context.TODO(), current.ID, imgRmForm.Force, imgRmForm.PruneChildren); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond()
	}

	ctx.Respond(msgs.ImageRemovedSuccessfully)

	images := h.updateImagesList(userID)
	if len(images) == 0 {
		return ctx.Send("there is no image")
	}
	current = images[0]

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImagesList(0),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ImageRemoveForce(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.InvalidButton)
	}

	current := session.GetImages()[index]
	imgRmForm := session.GetImageRemoveForm()
	imgRmForm.Force = !imgRmForm.Force
	session.SetImageRemoveForm(imgRmForm.Force, imgRmForm.PruneChildren)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(index, imgRmForm.Force, imgRmForm.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ImageRemovePruneChildren(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.InvalidButton)
	}

	current := session.GetImages()[index]
	imgRmForm := session.GetImageRemoveForm()
	imgRmForm.PruneChildren = !imgRmForm.PruneChildren
	session.SetImageRemoveForm(imgRmForm.Force, imgRmForm.PruneChildren)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(index, imgRmForm.Force, imgRmForm.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}
