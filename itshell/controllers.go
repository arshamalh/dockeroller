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

func StageGates() int {
	fmt.Print(msgs.Gates)
	if stage := getStage(); stage != 0 {
		return stage + 10
	}
	return 0
}

func StageTelegram() int {
	fmt.Print(msgs.Telegram)
	var token string
	getInput("Token: ", &token)
	var username string
	getInput("Username: @", &username)
	fmt.Println("Username and token successfully sat!")
	time.Sleep(time.Second * 3)
	return 0
}

func StageAPI() int {
	fmt.Print(msgs.API)
	var port int
	getInput("Port: ", &port)
	var password string
	getInput("Password: ", &password)
	fmt.Println("Port and password successfully sat!")
	time.Sleep(time.Second * 3)
	return 0
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
