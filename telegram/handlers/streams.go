package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/tbb"
)

// logLines is how many of the latest log lines the logs message shows.
const logLines = 20

// A stream keeps a message updated with live logs or stats. Each chat has at
// most one. It runs in its own goroutine, outside the handler: tbb processes
// one update per chat at a time, so a stream running inside a handler would
// block that chat's Back button.
type stream struct {
	cancel context.CancelFunc
	done   chan struct{}
}

// startStream replaces the chat's stream with run.
func (h *Handler) startStream(chatID int64, run func(ctx context.Context)) {
	h.stopStream(chatID)
	ctx, cancel := context.WithCancel(context.Background())
	s := &stream{cancel: cancel, done: make(chan struct{})}
	h.mu.Lock()
	h.streams[chatID] = s
	h.mu.Unlock()
	go func() {
		defer close(s.done)
		defer func() {
			h.mu.Lock()
			if h.streams[chatID] == s {
				delete(h.streams, chatID)
			}
			h.mu.Unlock()
		}()
		run(ctx)
	}()
}

// stopStream cancels the chat's stream and waits for it to finish, so it
// can't edit the message after the caller has redrawn it.
func (h *Handler) stopStream(chatID int64) {
	h.mu.Lock()
	s := h.streams[chatID]
	delete(h.streams, chatID)
	h.mu.Unlock()
	if s == nil {
		return
	}
	s.cancel()
	select {
	case <-s.done:
	case <-time.After(5 * time.Second):
	}
}

// streamTarget is the panel message a stream keeps editing.
func streamTarget(c Ctx) tbb.StoredMessage {
	return tbb.StoredMessage{
		MessageID: strconv.Itoa(c.Session().GetPreviousMessageWithKeyboard()),
		ChatID:    c.ID(),
	}
}

// editor edits the stream's message, skipping identical content (Telegram
// rejects those edits) and logging other failures.
type editor struct {
	bot  tbb.API
	msg  tbb.StoredMessage
	kb   *tbb.ReplyMarkup
	last string
	log  func(msg string, args ...any)
}

func (e *editor) show(text string) {
	if text == e.last {
		return
	}
	_, err := e.bot.Edit(e.msg, text, e.kb, tbb.ModeMarkdownV2)
	if err != nil && !errors.Is(err, tbb.ErrMessageNotModified) && !errors.Is(err, tbb.ErrSameMessageContent) {
		e.log("stream edit failed", "err", err)
		return
	}
	e.last = text
}

// ContainerLogs follows the container's logs, showing the latest lines.
func (h *Handler) ContainerLogs(c Ctx) error {
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	kb := keyboards.ContainerBack(c, index)
	_ = c.Respond()
	if err := c.CleanEdit(c.GetL8d(msgs.StreamStarting, "name", cont.Name), kb); err != nil {
		return err
	}
	ed := &editor{bot: c.Bot(), msg: streamTarget(c), kb: kb, log: c.Logger().Warn}
	name, id := cont.Name, cont.ID

	h.startStream(c.ID(), func(ctx context.Context) {
		rc, err := h.docker.ContainerLogs(ctx, id)
		if err != nil {
			ed.show(tbb.EscapeMarkdownV2("⚠️ " + err.Error()))
			return
		}
		defer rc.Close()

		lines := make(chan string, 256)
		go func() {
			defer close(lines)
			sc := bufio.NewScanner(rc)
			sc.Buffer(make([]byte, 64*1024), 1024*1024)
			for sc.Scan() {
				select {
				case lines <- sc.Text():
				case <-ctx.Done():
					return
				}
			}
		}()

		var buf []string
		ended := false
		tick := time.NewTicker(h.streamInterval)
		defer tick.Stop()
		ed.show(msgs.Logs(name, buf))
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-lines:
				if !ok {
					// The container stopped, or the daemon closed the stream.
					lines, ended = nil, true
					continue
				}
				buf = append(buf, line)
				if len(buf) > logLines {
					buf = buf[len(buf)-logLines:]
				}
			case <-tick.C:
				ed.show(msgs.Logs(name, buf))
				if ended {
					return
				}
			}
		}
	})
	return nil
}

// ContainerStats shows live CPU and memory usage.
func (h *Handler) ContainerStats(c Ctx) error {
	index, cont, ok, err := h.cachedContainer(c)
	if !ok {
		return err
	}
	kb := keyboards.ContainerBack(c, index)
	_ = c.Respond()
	if err := c.CleanEdit(c.GetL8d(msgs.StreamStarting, "name", cont.Name), kb); err != nil {
		return err
	}
	ed := &editor{bot: c.Bot(), msg: streamTarget(c), kb: kb, log: c.Logger().Warn}
	name, id := cont.Name, cont.ID

	h.startStream(c.ID(), func(ctx context.Context) {
		rc, err := h.docker.ContainerStats(ctx, id)
		if err != nil {
			ed.show(tbb.EscapeMarkdownV2("⚠️ " + err.Error()))
			return
		}
		defer rc.Close()

		samples := make(chan entities.Stats, 1)
		go func() {
			defer close(samples)
			dec := json.NewDecoder(rc)
			for {
				var s entities.Stats
				if err := dec.Decode(&s); err != nil {
					if !errors.Is(err, io.EOF) && ctx.Err() == nil {
						ed.log("decode stats", "err", err)
					}
					return
				}
				select {
				case samples <- s:
				case <-ctx.Done():
					return
				}
			}
		}()

		var latest *entities.Stats
		tick := time.NewTicker(h.streamInterval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case s, ok := <-samples:
				if !ok {
					if latest != nil {
						ed.show(msgs.Stats(name, *latest))
					}
					return
				}
				latest = &s
			case <-tick.C:
				if latest != nil {
					ed.show(msgs.Stats(name, *latest))
				}
			}
		}
	})
	return nil
}
