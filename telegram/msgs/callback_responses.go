package msgs

import "gopkg.in/telebot.v3"

var (
	CannotStartTheContainer           = NewCallbackResponse("We cannot start the container!")
	CannotStopTheContainer            = NewCallbackResponse("We cannot stop the container!")
	FillTheFormAndPressDone           = NewCallbackResponse("Please fill the form and press done")
	UnableToRemoveContainer           = NewCallbackResponse("Unable to remove container")
	ContainerRemovedSuccessfully      = NewCallbackResponse("Container removed successfully")
	InvalidButton                     = NewCallbackResponse("Invalid button 🤔️️️️️️")
	UnableToRemoveImage               = NewCallbackResponse("Unable to remove image")
	ImageRemovedSuccessfully          = NewCallbackResponse("Image removed successfully")
	NoContainer                       = NewCallbackResponse("There is either no containers or you should run /containers again!")
	StartedButUnavailableCurrentState = NewCallbackResponse("Container started, but we're not able to show current state.")
	StoppedButUnavailableCurrentState = NewCallbackResponse("Container stopped, but we're not able to show current state.")
)

func NewCallbackResponse(text string) *telebot.CallbackResponse {
	return &telebot.CallbackResponse{Text: text}
}
