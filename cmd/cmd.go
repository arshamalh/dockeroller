package cmd

import (
	"github.com/arshamalh/dockeroller/log"
	"github.com/spf13/cobra"
)

func Execute() {
	var root = &cobra.Command{
		Use:   "dockeroller",
		Short: "ChatOps application for controlling docker daemon through messengers such as Telegram",
	}

	log.InitGlobalLogger("log.json")

	registerStart(root)
	if err := root.Execute(); err != nil {
		log.Gl.Fatal(err.Error())
	}
}
