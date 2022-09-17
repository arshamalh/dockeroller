package main

import (
	"github.com/arshamalh/dockeroller/docker"
	"github.com/arshamalh/dockeroller/itshell"
	"github.com/arshamalh/dockeroller/telegram"
	"github.com/arshamalh/dockeroller/tools"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	docker := docker.New()

	// apiSrv := api.New(docker)

	telegramSrv, _ := telegram.New(docker)

	// Read Yaml configurations and set them
	configs, _ := tools.LoadYamlConfig()
	telegramSrv.SetConfig(configs.TelegramInfo)

	// App controller
	shell := itshell.New(telegramSrv) // api, slack, discord, ...
	shell.Run()
}
