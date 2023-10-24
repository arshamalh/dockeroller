package telegram

import "gopkg.in/telebot.v3"

var Commands = []telebot.Command{
	{
		Text:        "containers",
		Description: "list all existing containers",
	},
	{
		Text:        "images",
		Description: "list all existing images",
	},
}
