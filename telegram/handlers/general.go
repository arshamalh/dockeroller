package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func StartHandler(ctx tele.Context) error {
	return ctx.Send("hi " + ctx.Message().Sender.FirstName + "\n" + fmt.Sprint(ctx.Chat().ID) + "\n/containers /images")
}
