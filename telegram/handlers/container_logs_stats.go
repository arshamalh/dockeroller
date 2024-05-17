package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func (h *handler) ContainerLogs(ctx telebot.Context) error {
	// TODO: Starting from the beginning might cause confusion in long stream of errors, we should have a navigate till to the end button.
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	containerID, index := tools.ExtractIndexAndID(ctx.Data())

	quit := make(chan struct{})
	session.SetQuitChan(quit)
	stream, err := h.docker.ContainerLogs(context.TODO(), containerID)
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.NoLogsAvailable)
	}

	streamer := bufio.NewScanner(stream)
	queue := entities.NewQueue()
	for streamer.Scan() {
		select {
		case <-quit:
			return ctx.Respond(msgs.FinishingTheLogsStream)
		default:
		}
		newMsg := streamer.Text()
		queue.Push(newMsg)
		if queue.Length > entities.LOGS_QUEUE_LEN {
			queue.Pop()
		}

		// Omitted error by purpose (the error is just about not modified message because of repetitive content)
		err := ctx.Edit(
			msgs.FmtMono(queue.String()),
			keyboards.ContainerBack(index),
		)
		if err != nil && !errors.Is(err, telebot.ErrMessageNotModified) {
			log.Gl.Error(err.Error())
		}
		time.Sleep(entities.LOGS_PULL_INTERVAL)
	}
	return ctx.Respond()
}

func (h *handler) ContainerStats(ctx telebot.Context) error {
	userID := ctx.Chat().ID
	session := h.session.Get(userID)
	containerID, index := tools.ExtractIndexAndID(ctx.Data())

	quit := make(chan struct{})
	session.SetQuitChan(quit)
	stream, err := h.docker.ContainerStats(context.TODO(), containerID)
	if err != nil {
		log.Gl.Error(err.Error())
		return ctx.Respond(msgs.NoStatsAvailable)
	}

	streamer := bufio.NewScanner(stream)
	latestMsg := ""

	for streamer.Scan() {
		select {
		case <-quit:
			log.Gl.Debug("end of streaming stats for user", zap.Int64("used_id", userID))
			return ctx.Respond(msgs.FinishingTheStatsStream)
		default:
		}
		stats := entities.Stats{}
		err := json.Unmarshal(streamer.Bytes(), &stats)
		if err != nil {
			log.Gl.Error(err.Error())
		}

		if newMsg := msgs.FmtStats(stats); newMsg != latestMsg {
			err := ctx.Edit(
				newMsg, keyboards.ContainerBack(index),
			)
			if err != nil && !errors.Is(err, telebot.ErrMessageNotModified) {
				log.Gl.Error(err.Error())
			}
			latestMsg = newMsg
		}
		time.Sleep(entities.STATES_PULL_INTERVAL)
	}
	return ctx.Respond()
}
