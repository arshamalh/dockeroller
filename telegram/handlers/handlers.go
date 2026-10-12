// Package handlers implements dockeroller's Telegram flows on top of tbb.
package handlers

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/arshamalh/dockeroller/docker"
	"github.com/arshamalh/dockeroller/telegram/btns"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/tbb"
)

// Handler holds the dependencies of the flows and the live log/stats
// streams.
type Handler struct {
	app    *tbb.App[State]
	docker docker.Docker

	// dockerTimeout bounds every daemon call made on behalf of a button.
	dockerTimeout time.Duration
	// streamInterval is how often a log/stats stream refreshes its message.
	streamInterval time.Duration

	mu      sync.Mutex
	streams map[int64]*stream // by chat ID
}

// Option tunes a Handler (mostly for tests).
type Option func(*Handler)

// WithStreamInterval sets how often log/stats messages are refreshed
// (default 1s; Telegram rate-limits edits, so don't go much lower).
func WithStreamInterval(d time.Duration) Option {
	return func(h *Handler) { h.streamInterval = d }
}

// Register wires every command, button and scene of dockeroller into app.
func Register(app *tbb.App[State], d docker.Docker, opts ...Option) *Handler {
	h := &Handler{
		app:            app,
		docker:         d,
		dockerTimeout:  30 * time.Second,
		streamInterval: time.Second,
		streams:        make(map[int64]*stream),
	}
	for _, opt := range opts {
		opt(h)
	}

	app.Command("/start", h.Start)
	app.Command("/containers", h.ContainersList)
	app.Command("/images", h.ImagesList)

	// Containers
	app.Callback(btns.Containers, h.ContainersList)
	app.Callback(btns.ContNext, h.ContainersNav)
	app.Callback(btns.ContPrev, h.ContainersNav)
	app.Callback(btns.ContBack, h.ContainerBack)
	app.Callback(btns.ContStart, h.ContainerStart)
	app.Callback(btns.ContStop, h.ContainerStop)
	app.Callback(btns.ContPause, h.ContainerPause)
	app.Callback(btns.ContUnpause, h.ContainerUnpause)
	app.Callback(btns.ContLogs, h.ContainerLogs)
	app.Callback(btns.ContStats, h.ContainerStats)
	app.Callback(btns.ContRename, h.ContainerRename)
	app.Callback(btns.ContRemoveForm, h.ContainerRemoveForm)
	app.Callback(btns.ContRemoveForce, h.ContainerRemoveForce)
	app.Callback(btns.ContRemoveVolume, h.ContainerRemoveVolumes)
	app.Callback(btns.ContRemoveDone, h.ContainerRemoveDone)
	app.Scene(SceneRenameContainer).On(QNewContainerName, h.ContainerRenameText)

	// Images
	app.Callback(btns.Images, h.ImagesList)
	app.Callback(btns.ImgNext, h.ImagesNav)
	app.Callback(btns.ImgPrev, h.ImagesNav)
	app.Callback(btns.ImgBack, h.ImageBack)
	app.Callback(btns.ImgTag, h.ImageTag)
	app.Callback(btns.ImgRmForm, h.ImageRemoveForm)
	app.Callback(btns.ImgRmForce, h.ImageRemoveForce)
	app.Callback(btns.ImgRmPruneCh, h.ImageRemovePruneChildren)
	app.Callback(btns.ImgRmDone, h.ImageRemoveDone)
	app.Scene(SceneTagImage).On(QNewImageTag, h.ImageTagText)

	app.DefaultText(h.Start)
	app.DefaultCallback(func(c Ctx) error {
		return c.RespondL8d(msgs.ToastOutdatedButton)
	})

	return h
}

// Close stops every running stream.
func (h *Handler) Close() {
	h.mu.Lock()
	ids := make([]int64, 0, len(h.streams))
	for id := range h.streams {
		ids = append(ids, id)
	}
	h.mu.Unlock()
	for _, id := range ids {
		h.stopStream(id)
	}
}

// Start shows the main menu.
func (h *Handler) Start(c Ctx) error {
	h.stopStream(c.ID())
	c.ResetScene()
	if c.Callback() != nil {
		_ = c.Respond()
	}
	return c.CleanSend(c.GetL8d(msgs.Welcome, "name", c.Sender().FirstName), keyboards.Welcome(c))
}

func (h *Handler) dockerCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), h.dockerTimeout)
}

// errInvalidIndex is returned when a button's payload isn't a number.
var errInvalidIndex = errors.New("invalid index")

// buttonIndex parses the index a navigation or action button carries.
func buttonIndex(c Ctx) (int, error) {
	args := c.Args()
	if len(args) == 0 {
		return 0, errInvalidIndex
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, errInvalidIndex
	}
	return i, nil
}

// respondErr shows a failed Docker action as a toast. Telegram caps toasts
// at 200 characters, and daemon errors can be longer.
func respondErr(c Ctx, err error) error {
	msg := err.Error()
	if len(msg) > 180 {
		msg = msg[:177] + "..."
	}
	return c.RespondL8d(msgs.ToastActionFailed, "error", msg)
}
