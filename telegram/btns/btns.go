// Package btns declares the unique keys of every inline button. The same key
// builds the button (keyboards) and routes its presses (handlers), so the two
// can't drift apart. Changing a value breaks buttons on messages already sent.
package btns

import "github.com/arshamalh/tbb"

const (
	// Main menu
	Containers tbb.BtnKey = "conts"
	Images     tbb.BtnKey = "imgs"

	// Containers
	ContPrev         tbb.BtnKey = "contPrev"
	ContNext         tbb.BtnKey = "contNext"
	ContBack         tbb.BtnKey = "contBack"
	ContLogs         tbb.BtnKey = "contLogs"
	ContStop         tbb.BtnKey = "contStop"
	ContStart        tbb.BtnKey = "contStart"
	ContStats        tbb.BtnKey = "contStats"
	ContRename       tbb.BtnKey = "contRename"
	ContRemoveForm   tbb.BtnKey = "contRmForm"
	ContRemoveForce  tbb.BtnKey = "contFcRm"
	ContRemoveVolume tbb.BtnKey = "contRmVol"
	ContRemoveDone   tbb.BtnKey = "contRmDone"
	ContPause        tbb.BtnKey = "contPause"
	ContUnpause      tbb.BtnKey = "contUnpause"

	// Images
	ImgPrev      tbb.BtnKey = "imgPrev"
	ImgNext      tbb.BtnKey = "imgNext"
	ImgBack      tbb.BtnKey = "imgBack"
	ImgRmForm    tbb.BtnKey = "imgRmForm"
	ImgRmForce   tbb.BtnKey = "imgRmFc"
	ImgRmPruneCh tbb.BtnKey = "imgPnCh" // prune children
	ImgRmDone    tbb.BtnKey = "imgRmDone"
	ImgTag       tbb.BtnKey = "imgTag"
)
