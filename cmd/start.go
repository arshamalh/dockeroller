package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/arshamalh/dockeroller/docker"
	"github.com/arshamalh/dockeroller/telegram"
	"github.com/arshamalh/tbb"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// DefaultDSN keeps dockeroller's data (chat sessions and the Telegram
// accounts that used the bot) in a SQLite file next to where it runs.
// Nothing leaves the host.
const DefaultDSN = "file:dockeroller.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"

type startOptions struct {
	token     string
	whitelist []int64
	dsn       string
	apiURL    string
	debug     bool
}

func registerStart(root *cobra.Command) {
	var opts startOptions

	cmd := &cobra.Command{
		Use:   "start",
		Short: "starting telegram bot",
		Long: `Start the Telegram bot.

Every flag can also come from the environment (or a .env file):
  TOKEN / TELEGRAM__TOKEN   bot token from @BotFather
  WHITELISTED_IDS           comma separated chat IDs allowed to use the bot
  DB__DSN                   database (SQLite file or postgres:// URL)
  TELEGRAM__URL             Bot API server (default api.telegram.org)
The Docker daemon is found through DOCKER_HOST & co., like the docker CLI.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.token, "token", "t", "", "input your telegram token")
	cmd.Flags().Int64SliceVarP(&opts.whitelist, "whitelisted-ids", "w", nil, "a comma separated list of ids which are allowed to use this bot")
	cmd.Flags().StringVar(&opts.dsn, "db", "", "database DSN: a SQLite file (default ./dockeroller.db) or a postgres:// URL")
	cmd.Flags().StringVar(&opts.apiURL, "api-url", "", "Telegram Bot API server URL")
	cmd.Flags().BoolVar(&opts.debug, "debug", false, "verbose logging")
	root.AddCommand(cmd)
}

// fromEnv fills the options the flags left empty.
func (o *startOptions) fromEnv() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	if o.token == "" {
		o.token = firstEnv("TELEGRAM__TOKEN", "TOKEN")
	}
	if o.dsn == "" {
		o.dsn = firstEnv("DB__DSN")
	}
	if o.dsn == "" {
		o.dsn = DefaultDSN
	}
	if o.apiURL == "" {
		o.apiURL = firstEnv("TELEGRAM__URL")
	}
	if len(o.whitelist) == 0 {
		ids, err := parseIDs(firstEnv("WHITELISTED_IDS"))
		if err != nil {
			return fmt.Errorf("WHITELISTED_IDS: %w", err)
		}
		o.whitelist = ids
	}
	if o.token == "" {
		return errors.New("telegram can't start because no token is provided (--token or TOKEN)")
	}
	return nil
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func parseIDs(s string) ([]int64, error) {
	var ids []int64
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%q is not a chat ID", part)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func start(ctx context.Context, opts startOptions) error {
	level := slog.LevelInfo
	if opts.debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	if err := opts.fromEnv(); err != nil {
		return err
	}
	if len(opts.whitelist) == 0 {
		logger.Warn("no whitelisted IDs: the bot will only reply with each sender's chat ID; restart with -w <id> to allow someone")
	}

	d, err := docker.New()
	if err != nil {
		return err
	}
	defer d.Close()
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := d.Ping(pingCtx); err != nil {
		logger.Warn("docker daemon is not reachable yet", "err", err)
	}
	cancel()

	bot, err := telegram.New(tbb.Config{
		Telegram:        tbb.TelegramConfig{Token: opts.token, URL: opts.apiURL, PollerTimeout: 10 * time.Second},
		DB:              tbb.DBConfig{DSN: opts.dsn},
		DefaultLanguage: "en",
	}, d, opts.whitelist, logger)
	if err != nil {
		return err
	}
	defer bot.Handlers.Close()

	if ctx == nil {
		ctx = context.Background()
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	logger.Info("dockeroller started", "whitelisted", opts.whitelist)
	return bot.App.Run(ctx)
}
