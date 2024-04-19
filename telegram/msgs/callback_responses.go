package msgs

import "gopkg.in/telebot.v3"

var (
	CannotStartTheContainer           = NewCallbackResponse("We cannot start the container!")
	CannotStopTheContainer            = NewCallbackResponse("We cannot stop the container!")
	UnableToRemoveImage               = NewCallbackResponse("Unable to remove image")
	UnableToRemoveContainer           = NewCallbackResponse("Unable to remove container")
	UnableToFetchContainers           = NewCallbackResponse("Unable to fetch containers")
	UnableToFetchImages               = NewCallbackResponse("Unable to fetch images")
	FillTheFormAndPressDone           = NewCallbackResponse("Please fill the form and press done")
	InvalidButton                     = NewCallbackResponse("Invalid button 🤔️️️️️️")
	ContainerRemovedSuccessfully      = NewCallbackResponse("Container removed successfully")
	ImageRemovedSuccessfully          = NewCallbackResponse("Image removed successfully")
	NoContainer                       = NewCallbackResponse("There is either no containers or you should run /containers again!")
	NoImages                          = NewCallbackResponse("There is either no images or you should run /images again!")
	StartedButUnavailableCurrentState = NewCallbackResponse("Container started, but we're not able to show current state.")
	StoppedButUnavailableCurrentState = NewCallbackResponse("Container stopped, but we're not able to show current state.")
)

func NewCallbackResponse(text string) *telebot.CallbackResponse {
	return &telebot.CallbackResponse{Text: text}
}
