package api

import (
	"github.com/arshamalh/dockeroller/contracts"
)

type API interface {
	Start()
	Stop()
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
