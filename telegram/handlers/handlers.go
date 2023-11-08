package handlers

import (
	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/repo"
	"github.com/arshamalh/dockeroller/telegram/btns"
	"gopkg.in/telebot.v3"
)

type handler struct {
	docker  contracts.Docker
	bot     *telebot.Bot
	session repo.Session
	// log
}

func Register(bot *telebot.Bot, docker contracts.Docker, session repo.Session) {
	h := &handler{
		docker:  docker,
		bot:     bot,
		session: session,
	}

	// Command handlers
	h.bot.Handle("/start", StartHandler)
	h.bot.Handle("/containers", h.ContainersList)
	h.bot.Handle("/images", h.ImagesList)

	// Button handlers
	// Containers
	h.bot.Handle(btns.Containers.Key(), h.ContainersList)

	h.bot.Handle(btns.ContNext.Key(), h.ContainersNavBtn)
	h.bot.Handle(btns.ContPrev.Key(), h.ContainersNavBtn)
	h.bot.Handle(btns.ContLogs.Key(), h.ContainerLogs)
	h.bot.Handle(btns.ContStats.Key(), h.ContainerStats)
	h.bot.Handle(btns.ContBack.Key(), h.ContainersBackBtn)
	h.bot.Handle(btns.ContStop.Key(), h.ContainerStop)
	h.bot.Handle(btns.ContStart.Key(), h.ContainerStart)

	/// Images
	h.bot.Handle(btns.Images.Key(), h.ImagesList)

	h.bot.Handle(btns.ImgBack.Key(), h.ImagesBackBtn)
	h.bot.Handle(btns.ImgPrev.Key(), h.ImagesNavBtn)
	h.bot.Handle(btns.ImgNext.Key(), h.ImagesNavBtn)

}
