package handlers

import (
	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
)

// ContainerRemoveForm doesn't remove anything: it shows the options (force,
// remove volumes), reset to false for every container, and waits for Done.
//
// The cached list is deliberately not refreshed while the form is open, so
// toggling an option can't swap the container under the user's finger.
func (h *Handler) ContainerRemoveForm(c Ctx) error {
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	c.State().ContainerRemove = entities.ContainerRemoveForm{}
	_ = c.RespondL8d(msgs.ToastFillFormPressDone)
	return h.showContainerRemoveForm(c, index, cont)
}

func (h *Handler) showContainerRemoveForm(c Ctx, index int, cont *entities.Container) error {
	text := c.GetL8d(msgs.RemoveFormHint) + "\n\n" + c.GetL8d(msgs.Container, msgs.ContainerArgs(cont)...)
	return c.CleanEdit(text, keyboards.ContainerRemove(c, index, c.State().ContainerRemove))
}

func (h *Handler) ContainerRemoveForce(c Ctx) error {
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	form := &c.State().ContainerRemove
	form.Force = !form.Force
	_ = c.Respond()
	return h.showContainerRemoveForm(c, index, cont)
}

func (h *Handler) ContainerRemoveVolumes(c Ctx) error {
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	form := &c.State().ContainerRemove
	form.RemoveVolumes = !form.RemoveVolumes
	_ = c.Respond()
	return h.showContainerRemoveForm(c, index, cont)
}

// ContainerRemoveDone removes the container with the chosen options, then
// shows the (refreshed) first container.
func (h *Handler) ContainerRemoveDone(c Ctx) error {
	_, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	form := c.State().ContainerRemove
	ctx, cancel := h.dockerCtx()
	defer cancel()
	if err := h.docker.ContainerRemove(ctx, cont.ID, &form); err != nil {
		c.Logger().Warn("remove failed", "container", cont.Name, "err", err)
		return respondErr(c, err)
	}
	c.State().ContainerRemove = entities.ContainerRemoveForm{}
	_ = c.RespondL8d(msgs.ToastContainerRemoved)

	containers, err := h.refreshContainers(c)
	if err != nil {
		return dockerUnavailable(c, err)
	}
	if len(containers) == 0 {
		return c.CleanEdit(c.GetL8d(msgs.NoContainers), keyboards.Welcome(c))
	}
	return showContainer(c, 0, containers[0], "")
}
