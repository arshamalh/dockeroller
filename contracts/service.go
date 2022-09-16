package contracts

import "github.com/arshamalh/dockeroller/models"

type Service interface {
	Start()
	Stop()
	Info() models.ServiceInfo
	SetConfig(*Config)
}

type Services map[string]Service

func (ss Services) Get(name string) Service {
	return ss[name]
}

type Config map[string]any
