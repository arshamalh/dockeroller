package handlers

import (
	"context"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"github.com/docker/docker/api/types/filters"
)

// refreshContainers lists the containers and caches them in the session.
func (h *Handler) refreshContainers(c Ctx) ([]*entities.Container, error) {
	ctx, cancel := h.dockerCtx()
	defer cancel()
	containers, err := h.docker.ContainersList(ctx, filters.Args{})
	if err != nil {
		return nil, err
	}
	c.State().Containers = containers
	return containers, nil
}

// showContainer renders the container at index as the chat's panel: edited
// in place when coming from a button, sent fresh after a command or text.
func showContainer(c Ctx, index int, cont *entities.Container, prefix string) error {
	text := c.GetL8d(msgs.Container, msgs.ContainerArgs(cont)...)
	if prefix != "" {
		text = prefix + "\n\n" + text
	}
	kb := keyboards.Container(c, index, cont)
	if c.Callback() != nil {
		return c.CleanEdit(text, kb)
	}
	return c.CleanSend(text, kb)
}

// noContainers tells the user the list is empty, as a toast for buttons and
// as a message for commands.
func noContainers(c Ctx) error {
	if c.Callback() != nil {
		return c.RespondL8d(msgs.ToastNoContainers)
	}
	return c.SendL8d(msgs.NoContainers)
}

func dockerUnavailable(c Ctx, err error) error {
	c.Logger().Error("docker call failed", "err", err)
	if c.Callback() != nil {
		return respondErr(c, err)
	}
	return c.Send(c.GetL8d(msgs.DockerUnavailable, "error", err.Error()))
}

// ContainersList shows the first container (command or menu button).
func (h *Handler) ContainersList(c Ctx) error {
	h.stopStream(c.ID())
	c.ResetScene()
	containers, err := h.refreshContainers(c)
	if err != nil {
		return dockerUnavailable(c, err)
	}
	if len(containers) == 0 {
		return noContainers(c)
	}
	if c.Callback() != nil {
		_ = c.Respond()
	}
	return showContainer(c, 0, containers[0], "")
}

// ContainersNav handles Prev and Next; their payload is the target index,
// wrapped around the list.
func (h *Handler) ContainersNav(c Ctx) error {
	index, err := buttonIndex(c)
	if err != nil {
		return c.RespondL8d(msgs.ToastInvalidButton)
	}
	containers, err := h.refreshContainers(c)
	if err != nil {
		return dockerUnavailable(c, err)
	}
	if len(containers) == 0 {
		return noContainers(c)
	}
	index = tools.Indexer(index, len(containers))
	_ = c.Respond()
	return showContainer(c, index, containers[index], "")
}

// cachedContainer returns the container a button points at. When the cached
// list no longer has that index (it changed, or the session is new), the
// list is refreshed and the first container is shown instead; ok is false
// then and the caller should stop.
func (h *Handler) cachedContainer(c Ctx) (index int, cont *entities.Container, ok bool, err error) {
	index, err = buttonIndex(c)
	if err != nil {
		return 0, nil, false, c.RespondL8d(msgs.ToastInvalidButton)
	}
	containers := c.State().Containers
	if index >= 0 && index < len(containers) {
		return index, containers[index], true, nil
	}
	fresh, err := h.refreshContainers(c)
	if err != nil {
		return 0, nil, false, dockerUnavailable(c, err)
	}
	if len(fresh) == 0 {
		return 0, nil, false, noContainers(c)
	}
	_ = c.RespondL8d(msgs.ToastListChanged)
	return 0, nil, false, showContainer(c, 0, fresh[0], "")
}

// ContainerBack returns to the container view, from logs, stats, the remove
// form or the rename prompt. It also stops a running stream and leaves the
// rename scene.
func (h *Handler) ContainerBack(c Ctx) error {
	h.stopStream(c.ID())
	c.ResetScene()
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	_ = c.Respond()
	return showContainer(c, index, cont, "")
}

// containerAction runs a lifecycle action, then re-reads the container so
// the view and its buttons reflect the new state.
func (h *Handler) containerAction(c Ctx, toast string, action func(ctx context.Context, id string) error) error {
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	ctx, cancel := h.dockerCtx()
	defer cancel()
	if err := action(ctx, cont.ID); err != nil {
		c.Logger().Warn("container action failed", "container", cont.Name, "err", err)
		return respondErr(c, err)
	}
	// Re-list just this container: unlike inspect, the list gives the
	// familiar status text ("Up 2 seconds", "Exited (0) …").
	if fresh, err := h.docker.ContainersList(ctx, filters.NewArgs(filters.Arg("id", cont.ID))); err == nil && len(fresh) == 1 {
		cont = fresh[0]
		c.State().Containers[index] = cont
	}
	_ = c.RespondL8d(toast)
	return showContainer(c, index, cont, "")
}

func (h *Handler) ContainerStart(c Ctx) error {
	return h.containerAction(c, msgs.ToastContainerStarted, h.docker.ContainerStart)
}

func (h *Handler) ContainerStop(c Ctx) error {
	return h.containerAction(c, msgs.ToastContainerStopped, h.docker.ContainerStop)
}

func (h *Handler) ContainerPause(c Ctx) error {
	return h.containerAction(c, msgs.ToastContainerPaused, h.docker.ContainerPause)
}

func (h *Handler) ContainerUnpause(c Ctx) error {
	return h.containerAction(c, msgs.ToastContainerUnpaused, h.docker.ContainerUnpause)
}
