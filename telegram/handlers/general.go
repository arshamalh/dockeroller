package handlers

import (
	"strings"

	"github.com/arshamalh/dockeroller/telegram/keyboards"
	tele "gopkg.in/telebot.v3"
)

const welcomeMessage = `
Hello {name}, 
welcome to your bot,
You can use dockeroller to manage your docker daemon through different Messengers
e.g. list your images or containers:
`

func StartHandler(ctx tele.Context) error {
	newWelcomeMsg := strings.Replace(welcomeMessage, "{name}", ctx.Message().Sender.FirstName, -1)
	return ctx.Send(newWelcomeMsg, keyboards.Welcome())
}
