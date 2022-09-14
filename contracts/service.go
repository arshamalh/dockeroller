package contracts

import "github.com/arshamalh/dockeroller/models"

type Service interface {
	Start()
	Stop()
	Info() models.ServiceInfo
}
