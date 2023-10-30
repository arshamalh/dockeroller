package keyboards

import (
	"fmt"

	"github.com/arshamalh/dockeroller/telegram/btnkeys"
	tele "gopkg.in/telebot.v3"
)

func ContainersList(index int, is_on bool) *tele.ReplyMarkup {
	keyboard := &tele.ReplyMarkup{}
	var start_stop tele.Btn
	if is_on {
		start_stop = keyboard.Data("Stop", "stop")
	} else {
		start_stop = keyboard.Data("Start", "start")
	}
	keyboard.Inline(
		tele.Row{
			keyboard.Data("⬅", btnkeys.ContBack.String(), fmt.Sprint(index-1)),
			keyboard.Data("➡", "next", fmt.Sprint(index+1)),
		},
		tele.Row{
			start_stop,
			keyboard.Data("Remove", "remove"),
			keyboard.Data("Rename", "rename"),
		},
		tele.Row{
			keyboard.Data("Logs", "logs", fmt.Sprint(index)),
			keyboard.Data("Stats", "stats", fmt.Sprint(index)),
		},
	)

	return keyboard
}

func Back(index int, is_on bool) *tele.ReplyMarkup {
	keyboard := &tele.ReplyMarkup{}
	var start_stop tele.Btn
	if is_on {
		start_stop = keyboard.Data("Stop", "stop")
	} else {
		start_stop = keyboard.Data("Start", "start")
	}
	keyboard.Inline(
		tele.Row{
			keyboard.Data("⬅", "back_containers", fmt.Sprint(index)),
			start_stop,
		},
	)

	return keyboard
}
