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
}

func New(docker contracts.Docker) API {
	return &api{docker}
}

func (api *api) Start() {}
func (api *api) Stop()  {}
