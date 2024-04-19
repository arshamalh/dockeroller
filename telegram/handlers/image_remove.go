package handlers

import (
	"context"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

// This handler doesn't remove the image,
// it's just a form for setting the options of removing a image in the later "Done" step.
func (h *handler) ImageRemoveForm(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	current, _ := session.GetCurrentImage()
	ctx.Respond(msgs.FillTheFormAndPressDone)
	irf := current.RemoveForm

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(current.ID, irf.Force, irf.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}

// Removes the image with specified options
func (h *handler) ImageRemoveDone(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	imageID := ctx.Data()
	current, _ := session.GetCurrentImage()
	irf := current.RemoveForm

	if current.ID[:entities.LEN_IMG_TRIM] != imageID {
		log.Gl.Info("Removing another image when current is different", zap.String("current", current.String()))
		return ctx.Respond(msgs.UnableToRemoveContainer)
	}

	if err := h.docker.ImageRemove(context.TODO(), imageID, irf.Force, irf.PruneChildren); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.UnableToRemoveImage)
	}

	ctx.Respond(msgs.ImageRemovedSuccessfully)

	return h.ImagesList(ctx)
}

func (h *handler) ImageRemoveForce(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	imageID := ctx.Data()
	session := h.session.Get(userID)

	current, _ := session.GetCurrentImage()
	imgRmForm := session.GetImageRemoveForm()
	imgRmForm.Force = !imgRmForm.Force
	session.SetImageRemoveForm(imgRmForm.Force, imgRmForm.PruneChildren)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(imageID, imgRmForm.Force, imgRmForm.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ImageRemovePruneChildren(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	current, _ := session.GetCurrentImage()
	imgRmForm := session.GetImageRemoveForm()
	imgRmForm.PruneChildren = !imgRmForm.PruneChildren
	session.SetImageRemoveForm(imgRmForm.Force, imgRmForm.PruneChildren)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(current.ID, imgRmForm.Force, imgRmForm.PruneChildren),
		telebot.ModeMarkdownV2,
	)
}
