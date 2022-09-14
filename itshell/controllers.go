package itshell

import (
	"fmt"
	"time"

	"github.com/arshamalh/dockeroller/itshell/msgs"
)

func StageWelcome() int {
	fmt.Print(msgs.Welcome)
	return getStage()
}

func StageHelp() int {
	fmt.Print(msgs.Help)
	return getStage()
}

func StageGates(telegram bool, api bool) int {
	fmt.Print(msgs.Gates(telegram, api))
	if stage := getStage(); stage != 0 {
		return stage + 10
	}
	return 0
}

func StageTelegram() (stage int, token string, username string) {
	fmt.Print(msgs.Telegram)
	getInput("Token: ", &token)
	getInput("Username: @", &username)
	fmt.Println("Username and token successfully sat!")
	time.Sleep(time.Second * 3)
	return
}

func StageAPI() (stage int, port int, password string) {
	fmt.Print(msgs.API)
	getInput("Port: ", &port)
	getInput("Password: ", &password)
	fmt.Println("Port and password successfully sat!")
	time.Sleep(time.Second * 3)
	return
}

func getStage() (stage int) {
	fmt.Print("> ")
	fmt.Scanln(&stage)
	return
}

// Helper function to print a message and get a value (by populating a pointer)
func getInput(print string, value interface{}) {
	fmt.Print(print)
	fmt.Scanln(value)
}
