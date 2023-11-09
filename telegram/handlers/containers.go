package handlers

import (
	"bufio"
	"encoding/json"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainersNavBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	conts := h.session.GetContainers(userID)
	if len(conts) == 0 {
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "There is either no containers or you should run /containers again!",
			},
		)
	}
	index = tools.Indexer(index, len(conts))
	current := conts[index]

	containerIsOn := current.State == models.ContainerStateRunning
	err = ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, containerIsOn),
		telebot.ModeMarkdownV2,
	)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	return ctx.Respond()
}

func (h *handler) ContainersBackBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	h.session.GetQuitChan(userID) <- struct{}{}
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]

	containerIsOn := current.State == models.ContainerStateRunning
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, containerIsOn),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainersList(ctx telebot.Context) error {
	ctx.Respond()
	userID := ctx.Chat().ID
	containers := h.docker.ContainersList()
	h.session.SetContainers(userID, containers)
	current := containers[0]

	containerIsOn := current.State == models.ContainerStateRunning
	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, containerIsOn),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerLogs(ctx telebot.Context) error {
	// TODO: Starting from the beginning might cause confusion in long stream of errors, we should have a navigate till to the end button.
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]
	quit := make(chan struct{})
	h.session.SetQuitChan(userID, quit)
	stream, err := h.docker.ContainerLogs(current.ID)
	if err != nil {
		log.Gl.Error(err.Error())
	}

	streamer := bufio.NewScanner(stream)
	queue := models.NewQueue()
	for streamer.Scan() {
		select {
		case <-quit:
			return nil
		default:
			newMsg := streamer.Text()
			queue.Push(newMsg)
			if queue.Length > 10 { // TODO: 10 should not be hard-coded
				queue.Pop()
			}

			// Omitted error by purpose (the error is just about not modified message because of repetitive content)
			ctx.Edit(
				queue.String(),
				keyboards.Back(index, true),
			)
			time.Sleep(time.Millisecond * 500)
			// TODO: sleeping time, not hardcoded, not too much, not so little (under 500 millisecond would be annoying)
		}
	}
	return ctx.Respond()
}

func (h *handler) ContainerStats(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]
	quit := make(chan struct{})
	h.session.SetQuitChan(userID, quit)
	stream, err := h.docker.ContainerStats(current.ID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	streamer := bufio.NewScanner(stream)
	latestMsg := ""
	for streamer.Scan() {
		select {
		case <-quit:
			log.Gl.Debug("end of streaming stats for user", zap.Int64("used_id", userID))
			return nil
		default:
			stats := models.Stats{}
			err := json.Unmarshal(streamer.Bytes(), &stats)
			if err != nil {
				log.Gl.Error(err.Error())
			}

			if newMsg := msgs.FmtStats(stats); newMsg != latestMsg {
				err := ctx.Edit(
					newMsg,
					keyboards.Back(index, true),
					telebot.ModeMarkdownV2,
				)
				if err != nil {
					log.Gl.Error(err.Error())
				}
				latestMsg = newMsg
			}
			time.Sleep(time.Second)
		}
	}
	return ctx.Respond()
}

func (h *handler) ContainerStart(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]
	if err := h.docker.ContainerStart(current.ID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "We cannot start the container!",
			},
		)
	}

	current, err = h.docker.GetContainer(current.ID)
	if err != nil {
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "Container started, but we're not able to show current state.",
			},
		)
	}
	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, true),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerStop(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]
	if err := h.docker.ContainerStop(current.ID); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "We cannot stop the container!",
			},
		)
	}

	current, err = h.docker.GetContainer(current.ID)
	if err != nil {
		return ctx.Respond(
			&telebot.CallbackResponse{
				Text: "Container stopped, but we're not able to show current state.",
			},
		)
	}

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, false),
		telebot.ModeMarkdownV2,
	)
}

// This handler doesn't remove the container,
// it's just a form for setting the options of removing a container in the later "Done" step.
func (h *handler) ContainerRemoveForm(ctx telebot.Context) error {
	ctx.Respond(&telebot.CallbackResponse{Text: "Please fill the form and press done"})
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]
	cRemoveForm := h.session.GetContainerRemoveForm(userID)
	if cRemoveForm == nil {
		cRemoveForm = h.session.SetContainerRemoveForm(userID, false, false)
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
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]
	cRemoveForm := h.session.GetContainerRemoveForm(userID)

	if err := h.docker.ContainerRemove(current.ID, cRemoveForm); err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{Text: "Unable to remove container"})
	}

	ctx.Respond(&telebot.CallbackResponse{Text: "Container removed successfully"})

	containers := h.docker.ContainersList()
	h.session.SetContainers(userID, containers)
	current = containers[0]
	containerIsOn := current.State == models.ContainerStateRunning

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, containerIsOn),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerRemoveForce(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{
			Text: "Invalid button 🤔️️️️️️",
		})
	}

	current := h.session.GetContainers(userID)[index]
	cRemoveForm := h.session.GetContainerRemoveForm(userID)
	cRemoveForm.Force = !cRemoveForm.Force
	h.session.SetContainerRemoveForm(userID, cRemoveForm.Force, cRemoveForm.RemoveVolumes)

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(index, cRemoveForm.Force, cRemoveForm.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainerRemoveVolumes(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(&telebot.CallbackResponse{
			Text: "Invalid button 🤔️️️️️️",
		})
	}

	current := h.session.GetContainers(userID)[index]
	cRemoveForm := h.session.GetContainerRemoveForm(userID)
	cRemoveForm.RemoveVolumes = !cRemoveForm.RemoveVolumes
	h.session.SetContainerRemoveForm(userID, cRemoveForm.Force, cRemoveForm.RemoveVolumes)

	return ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainerRemove(index, cRemoveForm.Force, cRemoveForm.RemoveVolumes),
		telebot.ModeMarkdownV2,
	)
}
