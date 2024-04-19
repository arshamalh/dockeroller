package handlers

import (
	"context"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/docker/docker/api/types/filters"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

// This handler doesn't remove the container,
// it's just a form for setting the options of removing a container in the later "Done" step.
func (h *handler) ContainerRemoveForm(ctx telebot.Context) error {
	ctx.Respond(msgs.FillTheFormAndPressDone)
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	containerID := ctx.Data()

	current, index := session.GetCurrentContainer()
	crf := current.RemoveForm

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(containerID, index, crf.Force, crf.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}

// Removes the container with specified options
func (h *handler) ContainerRemoveDone(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	containerID := ctx.Data()
	session := h.session.Get(userID)

	current, _ := session.GetCurrentContainer()
	if current.ID[:entities.LEN_CONT_TRIM] != containerID {
		log.Gl.Info("Removing another container when current is different", zap.String("current", current.String()))
		return ctx.Respond(msgs.UnableToRemoveContainer)
	}

	if err := h.docker.ContainerRemove(context.TODO(), containerID, current.RemoveForm); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.UnableToRemoveContainer)
	}

	ctx.Respond(msgs.ContainerRemovedSuccessfully)

	containers, err := h.docker.ContainersList(context.TODO(), filters.Args{})
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Send("can't fetch containers") // TODO: Inform user we can't fetch containers
	}

	if len(containers) == 0 {
		// TODO: if container removed and there is no container,
		// get back to the start menu, no extra action needed.
		return ctx.Send("there is no container")
	}

	current = containers[0]
	session.SetCurrentContainer(current, 0)
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(current.ID, 0, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerRemoveForce(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	current, index := session.GetCurrentContainer()
	crf := current.RemoveForm
	crf.Force = !crf.Force
	session.SetCurrentContainer(current, index)

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(current.ID, index, crf.Force, crf.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerRemoveVolumes(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	current, index := session.GetCurrentContainer()
	crf := current.RemoveForm
	crf.RemoveVolumes = !crf.RemoveVolumes
	session.SetCurrentContainer(current, index)

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(current.ID, index, crf.Force, crf.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}
