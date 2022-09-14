package api

import (
	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
)

type API interface {
	Start()
	Stop()
	Info() models.ServiceInfo
}

type api struct {
	docker contracts.Docker
	isOn   bool
}

func New(docker contracts.Docker) API {
	return &api{docker, false}
}

func (api *api) Start() {}
func (api *api) Stop()  {}
func (api api) Info() models.ServiceInfo {
	return models.ServiceInfo{
		Name: "api",
		IsOn: api.isOn,
	}
}
