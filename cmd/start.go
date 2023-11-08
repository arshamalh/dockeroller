package cmd

import (
	"os"
	"time"

	"github.com/arshamalh/dockeroller/contracts"
	"github.com/arshamalh/dockeroller/docker"
	"github.com/arshamalh/dockeroller/log"
	"github.com/arshamalh/dockeroller/repo/ephemeral"
	tpkg "github.com/arshamalh/dockeroller/telegram"
	"github.com/arshamalh/dockeroller/telegram/handlers"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func registerStart(root *cobra.Command) {
	var token string
	whitelistedIDS := make([]int64, 0)

	cmd := &cobra.Command{
		Use:   "start",
		Short: "starting telegram bot",
		Run: func(cmd *cobra.Command, args []string) {
			start(token, whitelistedIDS)
		},
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "input your telegram token")
	cmd.Flags().Int64SliceVarP(&whitelistedIDS, "whitelisted-ids", "w", whitelistedIDS, "a comma separated list of ids which are allowed to use this bot")
	root.AddCommand(cmd)
}

func start(token string, whitelistedIDs []int64) {
	log.Gl.Info("server has started")

	if err := godotenv.Load(); err != nil {
		log.Gl.Error(err.Error())
	}

	docker := docker.New()

	// apiSrv := api.New(docker)
	if token == "" {
		if os.Getenv("TOKEN") != "" {
			token = os.Getenv("TOKEN")
		} else {
			log.Gl.Error("telegram can't start because no token is provided")
		}
	}
	startTelegram(docker, token, whitelistedIDs)
}

func startTelegram(docker contracts.Docker, token string, whitelistedIDs []int64) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Gl.Error(err.Error())
	}
	session := ephemeral.New()
	handlers.Register(bot, docker, session)
	// Middlewares
	bot.Use(middleware.Whitelist(whitelistedIDs...))
	// TODO: Disabled logger middleware for now.
	// bot.Use(LoggerMiddleware(log.Gl))
	if err := bot.SetCommands(tpkg.Commands); err != nil {
		log.Gl.Error(err.Error())
	}
	bot.Start()
}
