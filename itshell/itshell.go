package itshell

import (
	"github.com/arshamalh/dockeroller/models"
)

type ItShell interface {
	Run()
}

type Service interface {
	Start()
	Stop()
	Info() models.ServiceInfo
}

type itShell struct {
	services []Service
}

func New(services ...Service) ItShell {
	return &itShell{
		services: services,
	}
}

// Starts a blocking process that actively interact with user in the terminal
func (s *itShell) Run() {
	startServices(s.services...)
	runInteractiveShell()
}

// Start services that are on by default.
func startServices(services ...Service) {
	for _, srv := range services {
		if srv.Info().IsOn {
			srv.Start()
		}
	}
}

func runInteractiveShell() {
	for {
		// prints.SayWelcome()
	}
}
