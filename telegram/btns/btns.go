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
	ContPrev  BtnKey = "contPrev"
	ContNext  BtnKey = "contNext"
	ContBack  BtnKey = "contBack"
	ContLogs  BtnKey = "contLogs"
	ContStats BtnKey = "contStats"

	// Images
	ImgPrev BtnKey = "imgPrev"
	ImgNext BtnKey = "imgNext"
	ImgBack BtnKey = "imgBack"
)
