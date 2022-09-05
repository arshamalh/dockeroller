package telegram

import (
	"github.com/arshamalh/dockeroller/contracts"
)

type Telegram interface {
	Start()
	Stop()
}

type telegram struct {
	docker contracts.Docker
}

func New(docker contracts.Docker) Telegram {
	return &telegram{docker}
}

func (t *telegram) Start() {}
func (t *telegram) Stop()  {}
