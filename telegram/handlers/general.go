package handlers

import (
	"strings"

	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	tele "gopkg.in/telebot.v3"
)

func StartHandler(ctx tele.Context) error {
	newWelcomeMsg := strings.Replace(msgs.WelcomeMessage, "{name}", ctx.Message().Sender.FirstName, -1)
	return ctx.Send(newWelcomeMsg, keyboards.Welcome())
}
