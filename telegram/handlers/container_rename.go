package handlers

import (
	"strings"

	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/docker/docker/api/types/filters"
)

// ContainerRename asks for a new name; the next text message answers it
// (ContainerRenameText). Back cancels.
func (h *Handler) ContainerRename(c Ctx) error {
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	c.Enter(SceneRenameContainer, QNewContainerName)
	c.State().CurrentContainer = cont
	_ = c.Respond()
	return c.CleanEdit(c.GetL8d(msgs.ContainerNewName, "name", cont.Name), keyboards.ContainerBack(c, index))
}

// ContainerRenameText renames the container to the text the user sent.
// On failure the scene stays active so they can send another name; on
// success it ends, so later messages are not taken as names.
func (h *Handler) ContainerRenameText(c Ctx) error {
	cont := c.State().CurrentContainer
	if cont == nil {
		c.ResetScene()
		return h.Start(c)
	}
	newName := strings.TrimSpace(c.Text())

	ctx, cancel := h.dockerCtx()
	defer cancel()
	if err := h.docker.ContainerRename(ctx, cont.ID, newName); err != nil {
		c.Logger().Warn("rename failed", "container", cont.Name, "err", err)
		return c.CleanSend(
			c.GetL8d(msgs.ContainerRenameErr, "name", cont.Name, "error", err.Error()),
			keyboards.ContainerBack(c, h.containerIndex(c, cont.ID)),
		)
	}

	c.ResetScene()
	c.State().CurrentContainer = nil
	containers, err := h.docker.ContainersList(ctx, filters.Args{})
	if err != nil {
		return dockerUnavailable(c, err)
	}
	c.State().Containers = containers
	index := 0
	for i, x := range containers {
		if x.ID == cont.ID {
			index = i
			break
		}
	}
	if len(containers) == 0 {
		return c.SendL8d(msgs.NoContainers)
	}
	done := c.GetL8d(msgs.ContainerRenamed, "old", cont.Name, "new", newName)
	return showContainer(c, index, containers[index], done)
}

// containerIndex finds a container in the cached list (0 if absent).
func (h *Handler) containerIndex(c Ctx, id string) int {
	for i, x := range c.State().Containers {
		if x.ID == id {
			return i
		}
	}
	return 0
}
