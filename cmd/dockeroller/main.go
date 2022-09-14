package main

import (
	"github.com/arshamalh/dockeroller/api"
	"github.com/arshamalh/dockeroller/docker"
	"github.com/arshamalh/dockeroller/itshell"
	"github.com/arshamalh/dockeroller/telegram"
)

func main() {
	// Third Parties
	docker := docker.New()

	apiSrv := api.New(docker)

	telegramSrv := telegram.New(docker)

	// App controller
	shell := itshell.New(telegramSrv, apiSrv) // slack, discord, ...
	shell.Run()
}
