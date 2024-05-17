package keyboards

import (
	"fmt"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/telegram/btns"
	"gopkg.in/telebot.v3"
)

func ContainersList(containerID string, index int, containerIsOn bool) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	containerID = containerID[:entities.LEN_CONT_TRIM]

	keyboard.Inline(
		telebot.Row{
			keyboard.Data("Prev ⬅", btns.ContPrev.String(), fmt.Sprint(index-1)),
			keyboard.Data("Next ➡", btns.ContNext.String(), fmt.Sprint(index+1)),
		},
		telebot.Row{
			switchBtn(keyboard, containerID, containerIsOn),
			keyboard.Data("Remove 🗑", btns.ContRemoveForm.String(), containerID, fmt.Sprint(index)),
			keyboard.Data("Rename ✏️", btns.ContRename.String(), containerID, fmt.Sprint(index)),
		},
		telebot.Row{
			keyboard.Data("Logs 🪵", btns.ContLogs.String(), containerID, fmt.Sprint(index)),
			keyboard.Data("Stats 📊", btns.ContStats.String(), containerID, fmt.Sprint(index)),
		},
	)

	return keyboard
}

func ContainerRemove(containerID string, index int, force, removeVolumes bool) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	containerID = containerID[:entities.LEN_CONT_TRIM]

	keyboard.Inline(
		telebot.Row{
			keyboard.Data(
				fmt.Sprintf("Force: %t", force),
				string(btns.ContRemoveForce),
				containerID,
				fmt.Sprint(index),
			),
			keyboard.Data(
				fmt.Sprintf("Remove Volumes: %t", removeVolumes),
				string(btns.ContRemoveVolume),
				containerID,
				fmt.Sprint(index),
			),
		},
		telebot.Row{
			keyboard.Data(
				"Done",
				string(btns.ContRemoveDone),
				containerID,
				fmt.Sprint(index),
			),
		},
		telebot.Row{
			keyboard.Data(
				"Back ⬅",
				btns.ContBack.String(),
				containerID,
				fmt.Sprint(index),
			),
		},
	)

	return keyboard
}

func ContainerBack(index int) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		telebot.Row{
			keyboard.Data("Back ⬅", btns.ContBack.String()),
		},
	)

	return keyboard
}

func switchBtn(keyboard *telebot.ReplyMarkup, containerID string, containerIsOn bool) telebot.Btn {
	if containerIsOn {
		return keyboard.Data("Stop 🛑", btns.ContStop.String(), containerID)
	} else {
		return keyboard.Data("Start 🏃", btns.ContStart.String(), containerID)
	}
}
