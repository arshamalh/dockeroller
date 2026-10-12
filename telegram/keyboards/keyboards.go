// Package keyboards builds dockeroller's inline keyboards. Every button is
// built from its btns key, and every navigation button carries the index of
// the item it applies to.
package keyboards

import (
	"fmt"
	"strconv"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/telegram/btns"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/tbb"
)

// Labeler returns button labels; handlers pass their context.
type Labeler interface {
	GetL8dBtn(key string, kv ...string) string
}

func idx(i int) string { return strconv.Itoa(i) }

func inline(rows ...tbb.Row) *tbb.ReplyMarkup {
	kb := &tbb.ReplyMarkup{}
	kb.Inline(rows...)
	return kb
}

// Welcome is the main menu.
func Welcome(l Labeler) *tbb.ReplyMarkup {
	return inline(tbb.Row{
		btns.Images.Build(l.GetL8dBtn(msgs.BtnImages)),
		btns.Containers.Build(l.GetL8dBtn(msgs.BtnContainers)),
	})
}

// Container shows the actions that make sense for the container's state:
// a running container can be stopped or paused, a paused one unpaused or
// stopped, anything else started.
func Container(l Labeler, index int, c *entities.Container) *tbb.ReplyMarkup {
	i := idx(index)
	var lifecycle tbb.Row
	switch c.State {
	case entities.ContainerStateRunning:
		lifecycle = tbb.Row{
			btns.ContStop.Build(l.GetL8dBtn(msgs.BtnStop), i),
			btns.ContPause.Build(l.GetL8dBtn(msgs.BtnPause), i),
		}
	case entities.ContainerStatePaused:
		lifecycle = tbb.Row{
			btns.ContUnpause.Build(l.GetL8dBtn(msgs.BtnUnpause), i),
			btns.ContStop.Build(l.GetL8dBtn(msgs.BtnStop), i),
		}
	default:
		lifecycle = tbb.Row{btns.ContStart.Build(l.GetL8dBtn(msgs.BtnStart), i)}
	}
	return inline(
		tbb.Row{
			btns.ContPrev.Build(l.GetL8dBtn(msgs.BtnPrev), idx(index-1)),
			btns.ContNext.Build(l.GetL8dBtn(msgs.BtnNext), idx(index+1)),
		},
		lifecycle,
		tbb.Row{
			btns.ContLogs.Build(l.GetL8dBtn(msgs.BtnLogs), i),
			btns.ContStats.Build(l.GetL8dBtn(msgs.BtnStats), i),
		},
		tbb.Row{
			btns.ContRename.Build(l.GetL8dBtn(msgs.BtnRename), i),
			btns.ContRemoveForm.Build(l.GetL8dBtn(msgs.BtnRemove), i),
		},
	)
}

// ContainerRemove is the remove form: two toggles, Done, and Back to cancel.
func ContainerRemove(l Labeler, index int, form entities.ContainerRemoveForm) *tbb.ReplyMarkup {
	i := idx(index)
	return inline(
		tbb.Row{
			btns.ContRemoveForce.Build(l.GetL8dBtn(msgs.BtnForce, "value", fmt.Sprint(form.Force)), i),
			btns.ContRemoveVolume.Build(l.GetL8dBtn(msgs.BtnRemoveVolumes, "value", fmt.Sprint(form.RemoveVolumes)), i),
		},
		tbb.Row{
			btns.ContBack.Build(l.GetL8dBtn(msgs.BtnBack), i),
			btns.ContRemoveDone.Build(l.GetL8dBtn(msgs.BtnDone), i),
		},
	)
}

// ContainerBack returns to the container at index (from logs, stats or a
// rename prompt).
func ContainerBack(l Labeler, index int) *tbb.ReplyMarkup {
	return inline(tbb.Row{btns.ContBack.Build(l.GetL8dBtn(msgs.BtnBack), idx(index))})
}

// Image shows the image actions.
func Image(l Labeler, index int) *tbb.ReplyMarkup {
	i := idx(index)
	return inline(
		tbb.Row{
			btns.ImgPrev.Build(l.GetL8dBtn(msgs.BtnPrev), idx(index-1)),
			btns.ImgNext.Build(l.GetL8dBtn(msgs.BtnNext), idx(index+1)),
		},
		tbb.Row{
			btns.ImgTag.Build(l.GetL8dBtn(msgs.BtnTag), i),
			btns.ImgRmForm.Build(l.GetL8dBtn(msgs.BtnRemove), i),
		},
	)
}

// ImageRemove is the image remove form.
func ImageRemove(l Labeler, index int, form entities.ImageRemoveForm) *tbb.ReplyMarkup {
	i := idx(index)
	return inline(
		tbb.Row{
			btns.ImgRmForce.Build(l.GetL8dBtn(msgs.BtnForce, "value", fmt.Sprint(form.Force)), i),
			btns.ImgRmPruneCh.Build(l.GetL8dBtn(msgs.BtnPruneChildren, "value", fmt.Sprint(form.PruneChildren)), i),
		},
		tbb.Row{
			btns.ImgBack.Build(l.GetL8dBtn(msgs.BtnBack), i),
			btns.ImgRmDone.Build(l.GetL8dBtn(msgs.BtnDone), i),
		},
	)
}

// ImageBack returns to the image at index.
func ImageBack(l Labeler, index int) *tbb.ReplyMarkup {
	return inline(tbb.Row{btns.ImgBack.Build(l.GetL8dBtn(msgs.BtnBack), idx(index))})
}
