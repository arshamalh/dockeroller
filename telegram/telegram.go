package telegram

import (
	"time"

	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/pkg/session"
	"github.com/arshamalh/dockeroller/telegram/handlers"
	"github.com/arshamalh/dockeroller/tools"
	tele "gopkg.in/telebot.v3"
)

// Telegram interface and telegram struct are replacable in clean code architecture
// At the time of writing this comment, they all have common methods and fields.

type Telegram interface {
	Start()
	Stop()
	Info() models.ServiceInfo
	SetConfig(*contracts.Config)
}

type telegram struct {
	bot    *tele.Bot
	docker contracts.Docker
	isOn   bool
	config *contracts.Config
}

func New(docker contracts.Docker, config *contracts.Config) (*telegram, error) {
	token, err := tools.GetToken(config)
	if err != nil {
		return nil, err
	}
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}
	session := session.New()
	handlers.Register(bot, docker, session)
	bot.SetCommands(commands)
	return &telegram{bot, docker, true, config}, nil
}

func (t *telegram) Start() {
	t.isOn = true
	go t.bot.Start()
}
func (t *telegram) Stop() {
	t.isOn = false
	t.bot.Stop()
}

func (t telegram) Info() models.ServiceInfo {
	return models.ServiceInfo{
		Name: "telegram",
		IsOn: t.isOn,
	}
}
