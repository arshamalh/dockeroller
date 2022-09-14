package itshell

import "github.com/arshamalh/dockeroller/contracts"

type ItShell interface {
	Run()
}

type itShell struct {
	services []contracts.Service
}

func New(services ...contracts.Service) ItShell {
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
func startServices(services ...contracts.Service) {
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
