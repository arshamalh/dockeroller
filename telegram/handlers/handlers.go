package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
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

	// TODO, prev button here is specific for containers, should be defined for images again.
	// Button handlers
	h.bot.Handle("\fprev", h.PrevNextBtnHandler)
	h.bot.Handle("\fnext", h.PrevNextBtnHandler)
	h.bot.Handle("\fstats", h.StatsHandler)
	h.bot.Handle("\flogs", h.LogsHandler)
	h.bot.Handle("\fback_containers", h.BackContainersBtnHandler)
}

func StartHandler(ctx tele.Context) error {
	return ctx.Send("hi " + ctx.Message().Sender.FirstName + "\n" + fmt.Sprint(ctx.Chat().ID) + "\n/containers /images")
}

func (h *handler) PrevNextBtnHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	conts := h.session.Get(ctx.Chat().ID, "conts").([]*models.Container)
	index = tools.Indexer(index, len(conts))
	current := conts[index]
	err = ctx.Edit(
		msgs.FmtContainer(current),
		keyboards.ContainersList(index, false),
		tele.ModeMarkdownV2,
	)
	if err != nil {
		fmt.Println(err)
	}
	return ctx.Respond()
}

func (h *handler) BackContainersBtnHandler(ctx tele.Context) error {
	h.session.Get(ctx.Chat().ID, "quit_channel").(chan struct{}) <- struct{}{}
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := h.session.Get(ctx.Chat().ID, "conts").([]*models.Container)[index]
	return ctx.Edit(
		msgs.FmtContainer(current),
		// TODO: false and true passed for making keyboards are hardcoded but should be changed soon.
		keyboards.ContainersList(index, false),
		tele.ModeMarkdownV2,
	)
}

func (h *handler) ContainersHandler(ctx tele.Context) error {
	containers := h.docker.ContainersList()
	h.session.Set(ctx.Chat().ID, "conts", containers)
	current := containers[0]
	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, false),
		tele.ModeMarkdownV2,
	)
}

func (h *handler) ImagesHandler(ctx tele.Context) error {
	images := h.docker.ImagesList()
	h.session.Set(ctx.Chat().ID, "imgs", images)
	current := images[0]
	return ctx.Send(
		msgs.FmtImage(current),
		keyboards.ContainersList(0, false),
		tele.ModeMarkdownV2,
	)
}

func (h *handler) LogsHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := h.session.Get(ctx.Chat().ID, "conts").([]*models.Container)[index]
	quit := make(chan struct{})
	h.session.Set(ctx.Chat().ID, "quit_channel", quit)
	for i := 0; i < 10; i++ {
		select {
		case <-quit:
			return nil
		default:
			err := ctx.Edit(
				current.Name+" log "+fmt.Sprint(i),
				keyboards.Back(index, true),
			)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Second)
		}
	}
	return ctx.Respond()
}

func (h *handler) StatsHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := h.session.Get(ctx.Chat().ID, "conts").([]*models.Container)[index]
	quit := make(chan struct{})
	h.session.Set(ctx.Chat().ID, "quit_channel", quit)
	stream, err := h.docker.ContainerStats(current.ID)
	if err != nil {
		fmt.Println(err)
	}
	streamer := bufio.NewScanner(stream)
	latest_msg := ""
	for streamer.Scan() {
		select {
		case <-quit:
			return nil
		default:
			stats := models.Stats{}
			json.Unmarshal(streamer.Bytes(), &stats)
			msg := msgs.FmtStats(stats)
			if msg != latest_msg {
				err := ctx.Edit(
					msg,
					keyboards.Back(index, true),
					tele.ModeMarkdownV2,
				)
				if err != nil {
					fmt.Println(err)
				}
				latest_msg = msg
			}
			time.Sleep(time.Second)
		}
	}
	return ctx.Respond()
}
