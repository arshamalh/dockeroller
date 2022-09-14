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

// It will start a blocking shell session for controlling the gates and configs.
func runInteractiveShell() {
	var stage int = 0
	for {
		switch stage {
		case 0:
			stage = StageWelcome()
		case 1:
			stage = StageHelp()
		case 2:
			stage = StageGates()
		case 11:
			stage = StageTelegram()
		case 12:
			stage = StageAPI()
		}
	}
}
