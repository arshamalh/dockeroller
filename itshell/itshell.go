package itshell

import "github.com/arshamalh/dockeroller/contracts"

type ItShell interface {
	Run()
}

type itShell struct {
	services contracts.Services
}

func New(services ...contracts.Service) ItShell {
	srvs := make(contracts.Services)
	for _, srv := range services {
		srvs[srv.Info().Name] = srv
	}

	return &itShell{
		services: srvs,
	}
}

// Starts a blocking process that actively interact with user in the terminal
func (s *itShell) Run() {
	startServices(s.services)
	runInteractiveShell(s.services)
}

// Start services that are on by default.
func startServices(services contracts.Services) {
	for _, srv := range services {
		if srv.Info().IsOn {
			srv.Start()
		}
	}
}

// It will start a blocking shell session for controlling the gates and configs.
func runInteractiveShell(services contracts.Services) {
	var stage int = 0
	for {
		switch stage {
		case 0:
			stage = StageWelcome()
		case 1:
			stage = StageHelp()
		case 2:
			stage = StageGates(services.Get("telegram").Info().IsOn, services.Get("api").Info().IsOn)
		case 11:
			// TODO: Is it a good idea to read configs from interactive shell?
			//   So, We can make a yaml or env file out of entered configs.
			//   Then load it in another part of the app, or any other good idea for it?
			// var token, username string
			stage, _, _ = StageTelegram()
			// Set Config then Start
			telegramSrv := services.Get("telegram")
			// telegramSrv.SetConfig(&contracts.Config{
			// 	"token":    token,
			// 	"username": username,
			// })
			telegramSrv.Start()
		case 12:
			// Set Config then start
			stage, _, _ = StageAPI()
		}
	}
}
