package keyboards

import (
	"fmt"

	"github.com/arshamalh/dockeroller/telegram/btns"
	"gopkg.in/telebot.v3"
)

func ContainersList(index int, containerIsOn bool) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		telebot.Row{
			keyboard.Data("Prev ⬅", btns.ContPrev.String(), fmt.Sprint(index-1)),
			keyboard.Data("Next ➡", btns.ContNext.String(), fmt.Sprint(index+1)),
		},
		telebot.Row{
			switchBtn(keyboard, index, containerIsOn),
			keyboard.Data("Pause ⏸", btns.ContPause.String(), fmt.Sprint(index)),
			keyboard.Data("Remove 🗑", btns.ContRemoveForm.String(), fmt.Sprint(index)),
			keyboard.Data("Rename ✏️", btns.ContRename.String(), fmt.Sprint(index)),
		},
		telebot.Row{
			keyboard.Data("Logs 🪵", btns.ContLogs.String(), fmt.Sprint(index)),
			keyboard.Data("Stats 📊", btns.ContStats.String(), fmt.Sprint(index)),
		},
	)
	return keyboard
}

func ContainerRemove(index int, force, removeVolumes bool) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		telebot.Row{
			keyboard.Data(
				fmt.Sprintf("Force: %t", force),
				string(btns.ContRemoveForce),
				fmt.Sprint(index),
			),
			keyboard.Data(
				fmt.Sprintf("Remove Volumes: %t", removeVolumes),
				string(btns.ContRemoveVolume),
				fmt.Sprint(index),
			),
		},
		telebot.Row{
			keyboard.Data("Done", string(btns.ContRemoveDone), fmt.Sprint(index)),
		},
	)

	return keyboard
}

func ContainerBack(index int) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		telebot.Row{
			keyboard.Data("Back ⬅", btns.ContBack.String(), fmt.Sprint(index)),
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
