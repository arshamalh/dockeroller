package keyboards

import (
	"fmt"

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

func ContainersList(index int, containerIsOn bool) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		telebot.Row{
			keyboard.Data("⬅", btns.ContPrev.String(), fmt.Sprint(index-1)),
			keyboard.Data("➡", btns.ContNext.String(), fmt.Sprint(index+1)),
		},
		telebot.Row{
			switchBtn(keyboard, containerIsOn),
			keyboard.Data("Remove", btns.ContRemove.String()),
			keyboard.Data("Rename", btns.ContRename.String()),
		},
		telebot.Row{
			keyboard.Data("Logs", btns.ContLogs.String(), fmt.Sprint(index)),
			keyboard.Data("Stats", btns.ContStats.String(), fmt.Sprint(index)),
		},
	)

	return keyboard
}

func ImagesList(index int) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		telebot.Row{
			keyboard.Data("Prev ⬅", btns.ImgPrev.String(), fmt.Sprint(index-1)),
			keyboard.Data("Next ➡", btns.ImgNext.String(), fmt.Sprint(index+1)),
		},
		telebot.Row{
			keyboard.Data("Run 🏁", btns.ImgRun.String(), fmt.Sprint(index)),
		},
		telebot.Row{
			keyboard.Data("Remove 🗑", btns.ImgRemove.String(), fmt.Sprint(index)),
			keyboard.Data("Rename ✏️", btns.ImgRename.String(), fmt.Sprint(index)),
		},
	)

	return keyboard
}

func Back(index int, containerIsOn bool) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		telebot.Row{
			keyboard.Data("⬅", btns.ContBack.String(), fmt.Sprint(index)),
			switchBtn(keyboard, containerIsOn),
		},
	)

	return keyboard
}

func switchBtn(keyboard *telebot.ReplyMarkup, containerIsOn bool) telebot.Btn {
	if containerIsOn {
		return keyboard.Data("Stop", btns.ContStop.String())
	} else {
		return keyboard.Data("Start", btns.ContStart.String())
	}
}
