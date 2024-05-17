package middlewares

import (
	"github.com/arshamalh/dockeroller/log"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func LoggerMiddleware(zapLogger *zap.Logger) func(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			if err := next(ctx); err != nil {
				log.Gl.Error(err.Error())
				return err
			}
			return nil
		}
	}
}
