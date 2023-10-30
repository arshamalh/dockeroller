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
	h.bot.Handle("/containers", h.ContainersHandler)
	h.bot.Handle("/images", h.ImagesHandler)

	// Button handlers
	h.bot.Handle(btns.ContNext.Key(), h.PrevNextBtnHandler)
	h.bot.Handle(btns.ContPrev.Key(), h.PrevNextBtnHandler)
	h.bot.Handle(btns.ContLogs.Key(), h.LogsHandler)
	h.bot.Handle(btns.ContStats.Key(), h.StatsHandler)
	h.bot.Handle(btns.ContBack.Key(), h.BackContainersBtnHandler)
}
