package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"gopkg.in/telebot.v3"
)

func (h *handler) PrevNextBtnHandler(ctx telebot.Context) error {
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
		telebot.ModeMarkdownV2,
	)
	if err != nil {
		fmt.Println(err)
	}
	return ctx.Respond()
}

func (h *handler) BackContainersBtnHandler(ctx telebot.Context) error {
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
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainersHandler(ctx telebot.Context) error {
	containers := h.docker.ContainersList()
	h.session.Set(ctx.Chat().ID, "conts", containers)
	current := containers[0]
	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, false),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ImagesHandler(ctx telebot.Context) error {
	images := h.docker.ImagesList()
	h.session.Set(ctx.Chat().ID, "imgs", images)
	current := images[0]
	return ctx.Send(
		msgs.FmtImage(current),
		keyboards.ContainersList(0, false),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) LogsHandler(ctx telebot.Context) error {
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

func (h *handler) StatsHandler(ctx telebot.Context) error {
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
					telebot.ModeMarkdownV2,
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
