// Package telegram builds the dockeroller bot.
package telegram

import (
	"fmt"
	"log/slog"

	"github.com/arshamalh/dockeroller/docker"
	"github.com/arshamalh/dockeroller/telegram/handlers"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/tbb"
)

// Commands is the menu Telegram shows next to the message box.
var Commands = []tbb.Command{
	{Text: "containers", Description: "list all existing containers"},
	{Text: "images", Description: "list all existing images"},
	{Text: "start", Description: "show the main menu"},
}

// Bot is a configured, not yet started dockeroller bot.
type Bot struct {
	App      *tbb.App[handlers.State]
	Handlers *handlers.Handler
}

// Option customizes New.
type Option func(*options)

type options struct {
	app      []tbb.AppOption[handlers.State]
	handlers []handlers.Option
}

// WithAppOptions passes options to tbb.New (store, telebot settings, ...).
func WithAppOptions(opts ...tbb.AppOption[handlers.State]) Option {
	return func(o *options) { o.app = append(o.app, opts...) }
}

// WithHandlerOptions passes options to handlers.Register.
func WithHandlerOptions(opts ...handlers.Option) Option {
	return func(o *options) { o.handlers = append(o.handlers, opts...) }
}

// New builds the bot. Only chats in whitelist may use it.
func New(cfg tbb.Config, d docker.Docker, whitelist []int64, logger *slog.Logger, opts ...Option) (*Bot, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	appOpts := append([]tbb.AppOption[handlers.State]{tbb.WithLogger[handlers.State](logger)}, o.app...)
	app, err := tbb.New[handlers.State](cfg, appOpts...)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}
	app.I18n().Register("en", "", msgs.EN)
	app.Use(Whitelist(whitelist))

	h := handlers.Register(app, d, o.handlers...)
	if err := app.Bot().SetCommands(Commands); err != nil {
		logger.Warn("set bot commands", "err", err)
	}
	return &Bot{App: app, Handlers: h}, nil
}

// Whitelist only lets the given chats through. Others are told their chat ID,
// so the owner can whitelist it, and nothing else happens: the bot controls
// the Docker daemon, so the default is to deny.
func Whitelist(allowed []int64) tbb.Middleware[handlers.State] {
	set := make(map[int64]bool, len(allowed))
	for _, id := range allowed {
		set[id] = true
	}
	return func(next tbb.Handler[handlers.State]) tbb.Handler[handlers.State] {
		return func(c handlers.Ctx) error {
			if sender := c.Sender(); sender != nil && set[sender.ID] {
				return next(c)
			}
			c.Logger().Warn("unauthorized update", "sender", senderID(c))
			if c.Callback() != nil {
				return c.Respond()
			}
			if c.Message() == nil {
				return nil
			}
			return c.Send(c.GetL8d(msgs.Unauthorized, "id", fmt.Sprint(senderID(c))))
		}
	}
}

func senderID(c handlers.Ctx) int64 {
	if s := c.Sender(); s != nil {
		return s.ID
	}
	return c.ID()
}
