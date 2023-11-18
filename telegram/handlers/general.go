package handlers

import (
	"strings"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func StartHandler(ctx telebot.Context) error {
	newWelcomeMsg := strings.Replace(msgs.WelcomeMessage, "{name}", ctx.Message().Sender.FirstName, -1)
	return ctx.Send(newWelcomeMsg, keyboards.Welcome())
}

func (h *handler) General(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	scene := h.session.GetScene(userID)
	switch scene {
	case models.SceneRenameContainer:
		return h.ContainerRenameTextHandler(ctx)
	case models.SceneRenameImage:
		return h.ImageTagTextHandler(ctx)
	}

	log.Gl.Error("orphan scene", zap.Int64("user id", userID), zap.Int("current scene", int(scene)))
	return StartHandler(ctx)
}
