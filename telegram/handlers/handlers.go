package handlers

import (
	"github.com/arshamalh/dockeroller/docker"
	sessionPkg "github.com/arshamalh/dockeroller/session"
	"github.com/arshamalh/dockeroller/telegram/btns"
	"gopkg.in/telebot.v3"
)

type handler struct {
	docker  docker.Docker
	bot     *telebot.Bot
	session sessionPkg.Session
}

func NewHandler(bot *telebot.Bot, docker docker.Docker, session sessionPkg.Session) *handler {
	h := &handler{
		docker:  docker,
		bot:     bot,
		session: session,
	}

	// Command handlers
	h.bot.Handle("/start", h.StartHandler)
	h.bot.Handle("/containers", h.ContainersList)
	h.bot.Handle("/images", h.ImagesList)

	// Text Handler
	h.bot.Handle(telebot.OnText, h.General)

	// Button handlers
	// ** Containers ** //
	h.bot.Handle(btns.Containers.Key(), h.ContainersList)
	h.bot.Handle(btns.ContNext.Key(), h.ContainersNavBtn)
	h.bot.Handle(btns.ContPrev.Key(), h.ContainersNavBtn)
	h.bot.Handle(btns.ContLogs.Key(), h.ContainerLogs)
	h.bot.Handle(btns.ContStats.Key(), h.ContainerStats)
	h.bot.Handle(btns.ContBack.Key(), h.ContainersBackBtn)
	h.bot.Handle(btns.ContStop.Key(), h.ContainerStop)
	h.bot.Handle(btns.ContStart.Key(), h.ContainerStart)
	h.bot.Handle(btns.ContRemoveForm.Key(), h.ContainerRemoveForm)
	h.bot.Handle(btns.ContRemoveForce.Key(), h.ContainerRemoveForce)
	h.bot.Handle(btns.ContRemoveVolume.Key(), h.ContainerRemoveVolumes)
	h.bot.Handle(btns.ContRemoveDone.Key(), h.ContainerRemoveDone)
	h.bot.Handle(btns.ContRename.Key(), h.ContainerRename)

	// ** Images ** //
	h.bot.Handle(btns.Images.Key(), h.ImagesList)
	h.bot.Handle(btns.ImgBack.Key(), h.ImagesBackBtn)
	h.bot.Handle(btns.ImgPrev.Key(), h.ImagesNavBtn)
	h.bot.Handle(btns.ImgNext.Key(), h.ImagesNavBtn)
	h.bot.Handle(btns.ImgRmForm.Key(), h.ImageRemoveForm)
	h.bot.Handle(btns.ImgRmForce.Key(), h.ImageRemoveForce)
	h.bot.Handle(btns.ImgRmPruneCh.Key(), h.ImageRemovePruneChildren)
	h.bot.Handle(btns.ImgRmDone.Key(), h.ImageRemoveDone)
	h.bot.Handle(btns.ImgTag.Key(), h.ImageTag)

	return h
}
