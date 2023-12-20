package handlers

import (
	"context"
	"strconv"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/docker/docker/api/types/filters"
	"gopkg.in/telebot.v3"
)

// This handler doesn't remove the container,
// it's just a form for setting the options of removing a container in the later "Done" step.
func (h *handler) ContainerRemoveForm(ctx telebot.Context) error {
	ctx.Respond(&telebot.CallbackResponse{Text: "Please fill the form and press done"})
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetContainer(index)
	cRemoveForm := session.GetContainerRemoveForm()
	if cRemoveForm == nil {
		cRemoveForm = session.SetContainerRemoveForm(false, false)
	}

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(index, cRemoveForm.Force, cRemoveForm.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}

// Removes the container with specified options
func (h *handler) ContainerRemoveDone(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetContainer(index)
	cRemoveForm := session.GetContainerRemoveForm()

	if err := h.docker.ContainerRemove(current.ID, cRemoveForm); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{Text: "Unable to remove container"})
	}

	ctx.Respond(&telebot.CallbackResponse{Text: "Container removed successfully"})

	containers := h.docker.ContainersList(context.TODO(), filters.Args{})
	session.SetContainers(containers)
	if len(containers) == 0 {
		return ctx.Send("there is no container")
	}
	current = containers[0]

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, current.IsOn()),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerRemoveForce(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{
			Text: "Invalid button 🤔️️️️️️",
		})
	}

	current := session.GetContainer(index)
	cRemoveForm := session.GetContainerRemoveForm()
	cRemoveForm.Force = !cRemoveForm.Force
	session.SetContainerRemoveForm(cRemoveForm.Force, cRemoveForm.RemoveVolumes)

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(index, cRemoveForm.Force, cRemoveForm.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerRemoveVolumes(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{
			Text: "Invalid button 🤔️️️️️️",
		})
	}

	current := session.GetContainer(index)
	cRemoveForm := session.GetContainerRemoveForm()
	cRemoveForm.RemoveVolumes = !cRemoveForm.RemoveVolumes
	session.SetContainerRemoveForm(cRemoveForm.Force, cRemoveForm.RemoveVolumes)

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(index, cRemoveForm.Force, cRemoveForm.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}
