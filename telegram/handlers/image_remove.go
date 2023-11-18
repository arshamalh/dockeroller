package handlers

import (
	"context"
	"strconv"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"gopkg.in/telebot.v3"
)

// This handler doesn't remove the image,
// it's just a form for setting the options of removing a image in the later "Done" step.
func (h *handler) ImageRemoveForm(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{Text: "wrong button clicked!"})
	}
	// Informational Response
	ctx.Respond(&telebot.CallbackResponse{Text: "Please fill the form and press done"})

	current := h.session.GetImages(userID)[index]
	imgRmForm := h.session.GetImageRemoveForm(userID)
	if imgRmForm == nil {
		imgRmForm = h.session.SetImageRemoveForm(userID, false, false)
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
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetImages(userID)[index]
	imgRmForm := h.session.GetImageRemoveForm(userID)

	if err := h.docker.ImageRemove(context.TODO(), current.ID, imgRmForm.Force, imgRmForm.PruneChildren); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{Text: "Unable to remove image"})
	}

	ctx.Respond(&telebot.CallbackResponse{Text: "Image removed successfully"})

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
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{
			Text: "Invalid button 🤔️️️️️️",
		})
	}

	current := h.session.GetImages(userID)[index]
	imgRmForm := h.session.GetImageRemoveForm(userID)
	imgRmForm.Force = !imgRmForm.Force
	h.session.SetImageRemoveForm(userID, imgRmForm.Force, imgRmForm.PruneChildren)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(index, imgRmForm.Force, imgRmForm.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ImageRemovePruneChildren(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{
			Text: "Invalid button 🤔️️️️️️",
		})
	}

	current := h.session.GetImages(userID)[index]
	imgRmForm := h.session.GetImageRemoveForm(userID)
	imgRmForm.PruneChildren = !imgRmForm.PruneChildren
	h.session.SetImageRemoveForm(userID, imgRmForm.Force, imgRmForm.PruneChildren)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(index, imgRmForm.Force, imgRmForm.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) updateImagesList(userID int64) []*models.Image {
	images := h.docker.ImagesList()
	h.session.SetImages(userID, images)
	return images
}
