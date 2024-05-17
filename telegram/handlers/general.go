package handlers

import (
	"strings"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func (h *handler) StartHandler(ctx telebot.Context) error {
	newWelcomeMsg := strings.Replace(msgs.WelcomeMessage, "{name}", ctx.Message().Sender.FirstName, -1)
	return ctx.Send(msgs.FmtMono(newWelcomeMsg), keyboards.Welcome())
}

func (h *handler) General(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	scene := session.GetScene()
	switch scene {
	case entities.SceneRenameContainer:
		return h.ContainerRenameTextHandler(ctx)
	case entities.SceneRenameImage:
		return h.ImageTagTextHandler(ctx)
	}

	log.Gl.Error("orphan scene", zap.Int64("user id", userID), zap.Int("current scene", int(scene)))
	return h.StartHandler(ctx)
}

// We SHOULD respond to button click events,
// but when we don't want to,
// we should at least call the empty Respond so the button flash light (UI of waiting) would go away.
func (h *handler) EmptyResponder(ctx telebot.Context) {
	if err := ctx.Respond(); err != nil {
		log.Gl.Error(err.Error())
	}
}
