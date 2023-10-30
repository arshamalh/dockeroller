package handlers

import (
	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/telegram/btnkeys"
	tele "gopkg.in/telebot.v3"
)

type handler struct {
	docker  contracts.Docker
	bot     *tele.Bot
	session contracts.Session
	// log
}

func Register(bot *tele.Bot, docker contracts.Docker, session contracts.Session) {
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
	h.bot.Handle(btnkeys.ContNext.Key(), h.PrevNextBtnHandler)
	h.bot.Handle(btnkeys.ContPrev.Key(), h.PrevNextBtnHandler)
	h.bot.Handle(btnkeys.ContLogs.Key(), h.LogsHandler)
	h.bot.Handle(btnkeys.ContStats.Key(), h.StatsHandler)
	h.bot.Handle(btnkeys.ContBack.Key(), h.BackContainersBtnHandler)
}
