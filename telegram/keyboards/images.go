package keyboards

import (
	"fmt"

	"github.com/arshamalh/dockeroller/telegram/btns"
	"gopkg.in/telebot.v3"
)

func ImagesList(index int) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		telebot.Row{
			keyboard.Data("Prev ⬅", string(btns.ImgPrev), fmt.Sprint(index-1)),
			keyboard.Data("Next ➡", string(btns.ImgNext), fmt.Sprint(index+1)),
		},
		telebot.Row{
			keyboard.Data("Run 🏁", string(btns.ImgRun), fmt.Sprint(index)),
		},
		telebot.Row{
			keyboard.Data("Remove 🗑", string(btns.ImgRmForm), fmt.Sprint(index)),
			keyboard.Data("Rename ✏️", string(btns.ImgRename), fmt.Sprint(index)),
		},
	)

	return keyboard
}

func ImageRemove(index int, force, pruneChildren bool) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}

	keyboard.Inline(
		telebot.Row{
			keyboard.Data(
				fmt.Sprintf("Force: %t", force),
				string(btns.ImgRmForce),
				fmt.Sprint(index),
			),
			keyboard.Data(
				fmt.Sprintf("Prune Children: %t", pruneChildren),
				string(btns.ImgRmPruneCh),
				fmt.Sprint(index),
			),
		},
		telebot.Row{
			keyboard.Data("Done", string(btns.ImgRmDone), fmt.Sprint(index)),
		},
	)

	return keyboard
}
