package keyboards

import (
	"github.com/arshamalh/dockeroller/telegram/btns"
	"gopkg.in/telebot.v3"
)

func Welcome() *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		telebot.Row{
			keyboard.Data("Images 🌄", btns.Images.String()),
			keyboard.Data("Containers 📦", btns.Containers.String()),
		},
	)
	return keyboard
}
