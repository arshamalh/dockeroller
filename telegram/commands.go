package telegram

import tele "gopkg.in/telebot.v3"

var commands = []tele.Command{
	{
		Text:        "containers",
		Description: "list all existing containers",
	},
	{
		Text:        "images",
		Description: "list all existing images",
	},
}
