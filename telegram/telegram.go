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
	config *models.TelegramInfo
}

func New(docker contracts.Docker) *telegram {
	return &telegram{docker, false, nil}
}

func (t *telegram) Start() {}
func (t *telegram) Stop()  {}
func (t telegram) Info() models.ServiceInfo {
	return models.ServiceInfo{
		Name: "telegram",
		IsOn: t.isOn,
	}
}

func (t *telegram) SetConfig(config *models.TelegramInfo) {
	if config != nil {
		t.config = config
	}
}
