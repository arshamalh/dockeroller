package btns

type BtnKey string

// Turns button unique to a pair of (`\f` + button string value) so it can accessible for handlers.
func (bk BtnKey) Key() string {
	return "\f" + string(bk)
}

// Turns button unique directly to its string value without any extra information attached.
func (bk BtnKey) String() string {
	return string(bk)
}

const (
	// Containers
	ContPrev   BtnKey = "contPrev"
	ContNext   BtnKey = "contNext"
	ContBack   BtnKey = "contBack"
	ContLogs   BtnKey = "contLogs"
	ContStop   BtnKey = "contStop"
	ContStart  BtnKey = "contStart"
	ContStats  BtnKey = "contStats"
	ContRemove BtnKey = "contRemove"
	ContRename BtnKey = "contRename"

	// Images
	ImgPrev   BtnKey = "imgPrev"
	ImgNext   BtnKey = "imgNext"
	ImgBack   BtnKey = "imgBack"
	ImgRun    BtnKey = "imgRun"
	ImgRemove BtnKey = "imgRm"
	ImgRename BtnKey = "imgRename"
)
