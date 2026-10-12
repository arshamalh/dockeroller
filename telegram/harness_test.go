package telegram_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/internal/testkit"
	"github.com/arshamalh/dockeroller/telegram"
	"github.com/arshamalh/dockeroller/telegram/handlers"
	"github.com/arshamalh/tbb"
	"github.com/stretchr/testify/require"
)

// The user driving the tests, and their private chat.
var owner = &tbb.User{ID: 7, FirstName: "Ars_ham", Username: "arsham"}

// harness runs the real bot against a fake Telegram and a fake Docker. Every
// update is handled synchronously, so assertions can follow each action.
type harness struct {
	t      *testing.T
	tg     *testkit.Telegram
	docker *testkit.Docker
	store  tbb.Store
	bot    *telegram.Bot
	nextID int
}

func newHarness(t *testing.T, d *testkit.Docker, whitelist ...int64) *harness {
	t.Helper()
	if whitelist == nil {
		whitelist = []int64{owner.ID}
	}
	h := &harness{t: t, tg: testkit.NewTelegram(t), docker: d, store: tbb.NewMemoryStore(), nextID: 1000}
	h.start(whitelist)
	// Every test also proves that Telegram accepted every message.
	t.Cleanup(h.requireValidMessages)
	return h
}

// start (re)creates the bot on the same store, like a process restart.
func (h *harness) start(whitelist []int64) {
	h.t.Helper()
	if h.bot != nil {
		h.bot.Handlers.Close()
	}
	bot, err := telegram.New(
		tbb.Config{Telegram: tbb.TelegramConfig{Token: "123:test", URL: h.tg.URL()}},
		h.docker, whitelist, slog.New(slog.NewTextHandler(io.Discard, nil)),
		telegram.WithAppOptions(
			tbb.WithStore[handlers.State](h.store),
			tbb.WithSettings[handlers.State](func(s *tbb.Settings) {
				s.Offline = true
				s.Synchronous = true
			}),
		),
		telegram.WithHandlerOptions(handlers.WithStreamInterval(20*time.Millisecond)),
	)
	require.NoError(h.t, err)
	h.bot = bot
	h.t.Cleanup(bot.Handlers.Close)
}

func (h *harness) chat() *tbb.Chat {
	return &tbb.Chat{ID: owner.ID, Type: tbb.ChatPrivate, FirstName: owner.FirstName}
}

// send makes the user send a text message (or a /command).
func (h *harness) send(text string) {
	h.sendAs(owner, text)
}

func (h *harness) sendAs(user *tbb.User, text string) {
	h.nextID++
	h.bot.App.Bot().ProcessUpdate(tbb.Update{ID: h.nextID, Message: &tbb.Message{
		ID:     h.nextID,
		Text:   text,
		Sender: user,
		Chat:   &tbb.Chat{ID: user.ID, Type: tbb.ChatPrivate},
	}})
}

// click presses the inline button with this label on the newest message
// carrying it, and returns that message's ID.
func (h *harness) click(label string) int {
	h.t.Helper()
	msg, btn, ok := h.tg.FindButton(label)
	if !ok {
		h.t.Fatalf("no button %q; latest message: %v", label, h.tg.Latest())
	}
	h.press(msg.ID, btn.Data)
	return msg.ID
}

// press sends a callback with raw data, as if a button on message msgID was
// pressed.
func (h *harness) press(msgID int, data string) {
	h.nextID++
	h.bot.App.Bot().ProcessUpdate(tbb.Update{ID: h.nextID, Callback: &tbb.Callback{
		ID:      "cb",
		Sender:  owner,
		Data:    data,
		Message: &tbb.Message{ID: msgID, Chat: h.chat()},
	}})
}

// panel returns the message the bot keeps editing (the newest one).
func (h *harness) panel() testkit.Message {
	return h.tg.Latest()
}

// row returns the labels of the panel's i-th button row ("" row if absent,
// so a failing assertion reports instead of panicking).
func (h *harness) row(i int) []string {
	rows := h.panel().Labels()
	if i >= len(rows) {
		return nil
	}
	return rows[i]
}

// requireValidMessages fails if Telegram rejected any message, e.g. because
// of unescaped MarkdownV2.
func (h *harness) requireValidMessages() {
	h.t.Helper()
	for _, c := range h.tg.Rejected() {
		h.t.Errorf("Telegram rejected %s %q: %s", c.Method, c.Text(), c.Error)
	}
}

// Fixtures with names that are hostile to MarkdownV2.
func fixtures() *testkit.Docker {
	return testkit.NewDocker(
		[]*entities.Container{
			{ID: "c1", Name: "web-1", Image: "nginx:1.27", State: entities.ContainerStateRunning, Status: "Up 2 hours"},
			{ID: "c2", Name: "api_server", Image: "my/api:v2.1", State: entities.ContainerStatePaused, Status: "Up 1 hour (Paused)"},
			{ID: "c3", Name: "stopped.one", Image: "busybox", State: entities.ContainerStateExited, Status: "Exited (0) 3 days ago"},
		},
		[]*entities.Image{
			{ID: "sha256:863dd03b51905b838bbe0e32d14295bbe855412b2b23c74b4061551c1d071236", Size: 187_000_000,
				Tags: []string{"nginx:1.27", "nginx:latest"}, Status: entities.ImageStatusInUse, CreatedAt: "2026-10-01 10:00:00"},
			{ID: "sha256:d5318421c392e607be38a117e99b95999a490cb97ee2d4828591e48c0ddefae7", Size: 4_200,
				Status: entities.ImageStatusUnUsedDangling, CreatedAt: "2026-09-30 08:00:00"},
		},
	)
}

// endless returns a reader that writes a numbered line every few
// milliseconds until cancelled, like a busy container.
func endless(prefix string) func(ctx context.Context, id string) io.ReadCloser {
	return func(ctx context.Context, id string) io.ReadCloser {
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()
			for i := 1; ctx.Err() == nil; i++ {
				if _, err := fmt.Fprintf(pw, "%s %d\n", prefix, i); err != nil {
					return
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()
		return pr
	}
}

// lines returns a reader that yields the given log lines, then waits for
// cancellation (like a container that is still running).
func lines(ls ...string) func(ctx context.Context, id string) io.ReadCloser {
	return func(ctx context.Context, id string) io.ReadCloser {
		pr, pw := io.Pipe()
		go func() {
			for _, l := range ls {
				if _, err := pw.Write([]byte(l + "\n")); err != nil {
					return
				}
			}
			<-ctx.Done()
			pw.Close()
		}()
		return pr
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}
