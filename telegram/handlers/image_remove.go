package handlers

import (
	"context"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/docker/docker/api/types/filters"
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
	)
}

// Removes the image with specified options
func (h *handler) ImageRemoveDone(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	imageID := ctx.Data()
	current, _ := session.GetCurrentImage()
	irf := current.RemoveForm

	if current.ShortID() != imageID {
		log.Gl.Info(
			"Removing another image when current is different",
			zap.String("current in session", current.String()),
			zap.String("from callback", imageID),
		)
		return ctx.Respond(msgs.UnableToRemoveImage)
	}

	if err := h.docker.ImageRemove(context.TODO(), imageID, irf.Force, irf.PruneChildren); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.UnableToRemoveImage)
	}

	ctx.Respond(msgs.ImageRemovedSuccessfully)

	images, err := h.docker.ImagesList(context.TODO(), filters.Args{})
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Send("can't fetch containers") // TODO: Inform user we can't fetch containers
	}

	if len(images) == 0 {
		// TODO: if image removed and there is no image,
		// get back to the start menu, no extra action needed.
		return ctx.Send("there is no image")
	}

	current = images[0]
	session.SetCurrentImage(current, 0)
	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImagesList(0),
	)
}

func (h *handler) ImageRemoveForce(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	imageID := ctx.Data()
	session := h.session.Get(userID)
	current, index := session.GetCurrentImage()
	irf := current.RemoveForm
	irf.Force = !irf.Force
	session.SetCurrentImage(current, index)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(imageID, irf.Force, irf.PruneChildren),
	)
}

func (h *handler) ImageRemovePruneChildren(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	current, index := session.GetCurrentImage()
	irf := current.RemoveForm
	irf.PruneChildren = !irf.PruneChildren
	session.SetCurrentImage(current, index)

	return ctx.Edit(
		msgs.FmtImage(current),
		keyboards.ImageRemove(current.ID, irf.Force, irf.PruneChildren),
	)
}
