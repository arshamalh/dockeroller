package telegram

import (
	"os"
	"time"

	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/models"
	"github.com/arshamalh/dockeroller/pkg/session"
	"github.com/arshamalh/dockeroller/telegram/handlers"
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

func New(docker contracts.Docker) (*telegram, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}
	session := session.New()
	handlers.Register(bot, docker, session)
	bot.SetCommands(commands)
	return &telegram{bot, docker, true, nil}, nil
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

func (t *telegram) SetConfig(config *contracts.Config) {
	if config != nil {
		t.config = config
		t.bot.Token = (*config)["token"].(string)
	}
}
