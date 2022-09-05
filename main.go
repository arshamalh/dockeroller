package main

import (
	"github.com/arshamalh/dockeroller/repo/api"
	"github.com/arshamalh/dockeroller/repo/docker"
	"github.com/arshamalh/dockeroller/repo/itshell"
	"github.com/arshamalh/dockeroller/repo/telegram"
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
