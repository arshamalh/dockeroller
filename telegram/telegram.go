package telegram

import (
	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
)

// Telegram interface and telegram struct are replacable in clean code architecture
// At the time of writing this comment, they all have common methods and fields.

type Telegram interface {
	Start()
	Stop()
	Info() models.ServiceInfo
	SetConfig(*contracts.Config)
}

type telegram struct {
	docker contracts.Docker
	isOn   bool
	config *contracts.Config
}

func New(docker contracts.Docker) *telegram {
	return &telegram{docker, false, nil}
}

func (t *telegram) Start() {
	t.isOn = true
}
func (t *telegram) Stop() {
	t.isOn = false
}
func (t telegram) Info() models.ServiceInfo {
	return models.ServiceInfo{
		Name: "telegram",
		IsOn: t.isOn,
	}
}

func (t *telegram) SetConfig(config *contracts.Config) {
	if config != nil {
		t.config = config
	}
}
