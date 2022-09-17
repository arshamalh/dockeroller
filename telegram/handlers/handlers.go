package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type telegram struct {
	docker  contracts.Docker
	bot     *tele.Bot
	session contracts.Session
	// log
}

func Register(bot *tele.Bot, docker contracts.Docker, session contracts.Session) {
	t := &telegram{
		docker:  docker,
		bot:     bot,
		session: session,
	}
	// Middlewares
	telegram_superadmin_id, _ := strconv.ParseInt(os.Getenv("TELE_ADMIN"), 10, 64)
	t.bot.Use(middleware.Whitelist(telegram_superadmin_id))

	// Command handlers
	t.bot.Handle("/start", StartHandler)
	t.bot.Handle("/containers", t.ContainersHandler)
	t.bot.Handle("/images", t.ImagesHandler)

	// TODO, prev button here is specific for containers, should be defined for images again.
	// Button handlers
	t.bot.Handle("\fprev", t.PrevNextBtnHandler)
	t.bot.Handle("\fnext", t.PrevNextBtnHandler)
	t.bot.Handle("\fstats", t.StatsHandler)
	t.bot.Handle("\flogs", t.LogsHandler)
	t.bot.Handle("\fback_containers", t.BackContainersBtnHandler)
}

func StartHandler(ctx tele.Context) error {
	return ctx.Send("hi " + ctx.Message().Sender.FirstName + "\n" + fmt.Sprint(ctx.Chat().ID) + "\n/containers /images")
}

func (t *telegram) PrevNextBtnHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	conts := t.session.Get(ctx.Chat().ID, "conts").([]*models.Container)
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

func (t *telegram) BackContainersBtnHandler(ctx tele.Context) error {
	t.session.Get(ctx.Chat().ID, "quit_channel").(chan struct{}) <- struct{}{}
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := t.session.Get(ctx.Chat().ID, "conts").([]*models.Container)[index]
	return ctx.Edit(
		msgs.FmtContainer(current),
		// TODO: false and true passed for making keyboards are hardcoded but should be changed soon.
		keyboards.ContainersList(index, false),
		tele.ModeMarkdownV2,
	)
}

func (t *telegram) ContainersHandler(ctx tele.Context) error {
	containers := t.docker.ContainersList()
	t.session.Set(ctx.Chat().ID, "conts", containers)
	current := containers[0]
	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, false),
		tele.ModeMarkdownV2,
	)
}

func (t *telegram) ImagesHandler(ctx tele.Context) error {
	images := t.docker.ImagesList()
	t.session.Set(ctx.Chat().ID, "imgs", images)
	current := images[0]
	return ctx.Send(
		msgs.FmtImage(current),
		keyboards.ContainersList(0, false),
		tele.ModeMarkdownV2,
	)
}

func (t *telegram) LogsHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := t.session.Get(ctx.Chat().ID, "conts").([]*models.Container)[index]
	quit := make(chan struct{})
	t.session.Set(ctx.Chat().ID, "quit_channel", quit)
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

func (t *telegram) StatsHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := t.session.Get(ctx.Chat().ID, "conts").([]*models.Container)[index]
	quit := make(chan struct{})
	t.session.Set(ctx.Chat().ID, "quit_channel", quit)
	stream, err := t.docker.ContainerStats(current.ID)
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
