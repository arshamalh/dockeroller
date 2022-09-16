package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	tele "gopkg.in/telebot.v3"
)

func RegisterHandlers(bot *tele.Bot, docker contracts.Docker) {

	// Command handlers
	bot.Handle("/start", StartHandler)
	bot.Handle("/containers", ContainersHandler(docker))
	bot.Handle("/images", ImagesHandler(docker))

	// Button handlers
	bot.Handle("\fprev", PrevNextBtnHandler)
	bot.Handle("\fnext", PrevNextBtnHandler)
	bot.Handle("\fstats", StatsHandler)
	bot.Handle("\flogs", LogsHandler)
	bot.Handle("\fback_containers", BackContainersBtnHandler)
}

func StartHandler(ctx tele.Context) error {
	return ctx.Send("hi " + ctx.Message().Sender.FirstName + "\n/containers /images")
}

func PrevNextBtnHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := GetSession("conts").([]*models.Container)[index]
	ctx.Edit(
		msgs.FmtContainer(current),
		MakeContainerKeyboard(index, false),
		tele.ModeMarkdownV2,
	)
	return ctx.Respond()
}

func BackContainersBtnHandler(ctx tele.Context) error {
	GetSession("quit_channel").(chan struct{}) <- struct{}{}
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := GetSession("conts").([]*models.Container)[index]
	return ctx.Edit(
		msgs.FmtContainer(current),
		MakeContainerKeyboard(index, false),
		tele.ModeMarkdownV2,
	)
}

func ContainersHandler(docker contracts.Docker) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		containers := docker.ContainersList()
		SetSession("conts", containers)
		current := containers[0]
		return ctx.Send(
			msgs.FmtContainer(current),
			MakeContainerKeyboard(0, false),
			tele.ModeMarkdownV2,
		)
	}
}

func ImagesHandler(docker contracts.Docker) func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		images := docker.ImagesList()
		SetSession("imgs", images)
		current := images[0]
		return ctx.Send(
			msgs.FmtImage(current),
			MakeContainerKeyboard(0, false),
			tele.ModeMarkdownV2,
		)
	}
}

func LogsHandler(ctx tele.Context) error {
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := GetSession("conts").([]*models.Container)[index]
	quit := make(chan struct{})
	SetSession("quit_channel", quit)
	for i := 0; i < 10; i++ {
		select {
		case <-quit:
			return nil
		default:
			err := ctx.Edit(
				current.Name+" log "+fmt.Sprint(i),
				MakeBackKeyboard(index, true),
			)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Second)
		}
	}
	return ctx.Respond()
}

func StatsHandler(ctx tele.Context) error {
	return ctx.Respond()
}
