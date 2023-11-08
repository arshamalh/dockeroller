package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainersNavBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	conts := h.session.GetContainers(userID)
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

func (h *handler) ContainersBackBtn(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	h.session.GetQuitChan(userID) <- struct{}{}
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := h.session.GetContainers(userID)[index]
	return ctx.Edit(
		msgs.FmtContainer(current),
		// TODO: false and true passed for making keyboards are hardcoded but should be changed soon.
		keyboards.ContainersList(index, false),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) ContainersList(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	containers := h.docker.ContainersList()
	h.session.SetContainers(userID, containers)
	current := containers[0]
	return ctx.Send(
		msgs.FmtContainer(current),
		keyboards.ContainersList(0, false),
		telebot.ModeMarkdownV2,
	)
}

func (h *handler) Logs(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		fmt.Println(err)
	}
	current := h.session.GetContainers(userID)[index]
	quit := make(chan struct{})
	h.session.SetQuitChan(userID, quit)
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

func (h *handler) Stats(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := h.session.GetContainers(userID)[index]
	quit := make(chan struct{})
	h.session.SetQuitChan(userID, quit)
	stream, err := h.docker.ContainerStats(current.ID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	streamer := bufio.NewScanner(stream)
	latest_msg := ""
	for streamer.Scan() {
		select {
		case <-quit:
			log.Gl.Debug("end of streaming stats for user", zap.Int64("used_id", userID))
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
					log.Gl.Error(err.Error())
				}
				latest_msg = msg
			}
			time.Sleep(time.Second)
		}
	}
	return ctx.Respond()
}
