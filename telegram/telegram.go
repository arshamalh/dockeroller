package telegram

import (
	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
)

type Telegram interface {
	Start()
	Stop()
	Info() models.ServiceInfo
}

type telegram struct {
	docker contracts.Docker
	isOn   bool
}

func New(docker contracts.Docker) Telegram {
	return &telegram{docker, true}
}

func (t *telegram) Start() {}
func (t *telegram) Stop()  {}
func (t telegram) Info() models.ServiceInfo {
	return models.ServiceInfo{
		Name: "telegram",
		IsOn: t.isOn,
	}
}
