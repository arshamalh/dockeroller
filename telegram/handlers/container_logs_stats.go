package handlers

import (
	"bufio"
	"encoding/json"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainerLogs(ctx telebot.Context) error {
	// TODO: Starting from the beginning might cause confusion in long stream of errors, we should have a navigate till to the end button.
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	index, err := strconv.Atoi(ctx.Data())
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetContainer(index)
	quit := make(chan struct{})
	session.SetQuitChan(quit)
	stream, err := h.docker.ContainerLogs(current.ID)
	if err != nil {
		log.Gl.Error(err.Error())
	}

	streamer := bufio.NewScanner(stream)
	queue := entities.NewQueue()
	for streamer.Scan() {
		select {
		case <-quit:
			return nil
		default:
			newMsg := streamer.Text()
			queue.Push(newMsg)
			if queue.Length > 10 { // TODO: 10 should not be hard-coded
				queue.Pop()
			}

			// Omitted error by purpose (the error is just about not modified message because of repetitive content)
			ctx.Edit(
				queue.String(),
				keyboards.ContainerBack(index),
			)
			time.Sleep(time.Millisecond * 500)
			// TODO: sleeping time, not hardcoded, not too much, not so little (under 500 millisecond would be annoying)
		}
	}
	return ctx.Respond()
}

func (h *handler) ContainerStats(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	index, err := strconv.Atoi(ctx.Data())
	session := h.session.Get(userID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	current := session.GetContainer(index)
	quit := make(chan struct{})
	session.SetQuitChan(quit)
	stream, err := h.docker.ContainerStats(current.ID)
	if err != nil {
		log.Gl.Error(err.Error())
	}
	streamer := bufio.NewScanner(stream)
	latestMsg := ""
	for streamer.Scan() {
		select {
		case <-quit:
			log.Gl.Debug("end of streaming stats for user", zap.Int64("used_id", userID))
			return nil
		default:
			stats := entities.Stats{}
			err := json.Unmarshal(streamer.Bytes(), &stats)
			if err != nil {
				log.Gl.Error(err.Error())
			}

			if newMsg := msgs.FmtStats(stats); newMsg != latestMsg {
				err := ctx.Edit(
					newMsg,
					keyboards.ContainerBack(index),
					telebot.ModeMarkdownV2,
				)
				if err != nil {
					log.Gl.Error(err.Error())
				}
				latestMsg = newMsg
			}
			time.Sleep(time.Second)
		}
	}
	return ctx.Respond()
}
