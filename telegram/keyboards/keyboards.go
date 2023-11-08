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
			keyboard.Data("Prev ⬅", btns.ContPrev.String(), fmt.Sprint(index-1)),
			keyboard.Data("Next ➡", btns.ContNext.String(), fmt.Sprint(index+1)),
		},
		telebot.Row{
			switchBtn(keyboard, index, containerIsOn),
			keyboard.Data("Remove 🗑", btns.ContRemove.String(), fmt.Sprint(index)),
			keyboard.Data("Rename ✏️", btns.ContRename.String(), fmt.Sprint(index)),
		},
		telebot.Row{
			keyboard.Data("Logs 🪵", btns.ContLogs.String(), fmt.Sprint(index)),
			keyboard.Data("Stats 📊", btns.ContStats.String(), fmt.Sprint(index)),
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
			switchBtn(keyboard, index, containerIsOn),
		},
	)

	return keyboard
}

func switchBtn(keyboard *telebot.ReplyMarkup, index int, containerIsOn bool) telebot.Btn {
	if containerIsOn {
		return keyboard.Data("Stop 🛑", btns.ContStop.String(), fmt.Sprint(index))
	} else {
		return keyboard.Data("Start 🏃", btns.ContStart.String(), fmt.Sprint(index))
	}
}
