package middlewares

import (
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func LoggerMiddleware(zapLogger *zap.Logger) func(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			zapLogger.Info(ctx.Update().Message.Text)
			return next(ctx)
		}
	}
}
