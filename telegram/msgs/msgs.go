// Package msgs holds dockeroller's texts and the helpers that turn entities
// into template values.
//
// Templates are MarkdownV2 and are rendered with ctx.GetL8d, which keeps the
// template's own markup (`code` labels, *bold*) and fully escapes every
// value. Container names, image tags and user input therefore can't break
// Telegram's parser, which is what made the old images list fail on the
// "[" of its tags.
package msgs

import "github.com/arshamalh/tbb"

// Translation keys.
const (
	Welcome            = "welcome"
	Container          = "container"
	Image              = "image"
	ContainerRenamed   = "container.renamed"
	ContainerNewName   = "container.new_name"
	ContainerRenameErr = "container.rename_error"
	ImageTagged        = "image.tagged"
	ImageNewTag        = "image.new_tag"
	ImageTagErr        = "image.tag_error"
	RemoveFormHint     = "remove.form_hint"
	NoContainers       = "containers.none"
	NoImages           = "images.none"
	DockerUnavailable  = "docker.unavailable"
	StreamStarting     = "stream.starting"
	Unauthorized       = "unauthorized"

	// Callback toasts (plain text, no markup)
	ToastContainerStarted  = "toast.container_started"
	ToastContainerStopped  = "toast.container_stopped"
	ToastContainerPaused   = "toast.container_paused"
	ToastContainerUnpaused = "toast.container_unpaused"
	ToastContainerRemoved  = "toast.container_removed"
	ToastImageRemoved      = "toast.image_removed"
	ToastActionFailed      = "toast.action_failed"
	ToastListChanged       = "toast.list_changed"
	ToastInvalidButton     = "toast.invalid_button"
	ToastOutdatedButton    = "toast.outdated_button"
	ToastFillFormPressDone = "toast.fill_form"
	ToastNoContainers      = "toast.no_containers"
	ToastNoImages          = "toast.no_images"
)

// Button labels
const (
	BtnContainers    = "btn.containers"
	BtnImages        = "btn.images"
	BtnPrev          = "btn.prev"
	BtnNext          = "btn.next"
	BtnBack          = "btn.back"
	BtnStart         = "btn.start"
	BtnStop          = "btn.stop"
	BtnPause         = "btn.pause"
	BtnUnpause       = "btn.unpause"
	BtnRemove        = "btn.remove"
	BtnRename        = "btn.rename"
	BtnLogs          = "btn.logs"
	BtnStats         = "btn.stats"
	BtnTag           = "btn.tag"
	BtnDone          = "btn.done"
	BtnForce         = "btn.force"
	BtnRemoveVolumes = "btn.remove_volumes"
	BtnPruneChildren = "btn.prune_children"
)

// EN is the English text set.
var EN = tbb.Translations{
	tbb.LangNameKey: "English",

	Welcome: "Hello {name} 👋\n" +
		"Welcome to *Dockeroller*. You can manage your Docker daemon from here, " +
		"e.g. list your images or containers:",

	Container: "`Name:   ` {name}\n" +
		"`Image:  ` {image}\n" +
		"`State:  ` {state}\n" +
		"`Status: ` {status}",

	Image: "`ID:      ` `{id}`\n" +
		"`Size:    ` {size}\n" +
		"`Tags:    ` {tags}\n" +
		"`Status:  ` {status}\n" +
		"`Created: ` {created}",

	ContainerNewName:   "✏️ Send the new name for *{name}*:",
	ContainerRenamed:   "✅ *{old}* renamed to *{new}*",
	ContainerRenameErr: "⚠️ Couldn't rename *{name}*: {error}\nSend another name, or press Back.",
	ImageNewTag:        "🏷 Send the new tag for `{id}` (e.g. `myapp:1.2`):",
	ImageTagged:        "✅ `{id}` tagged *{tag}*",
	ImageTagErr:        "⚠️ Couldn't tag `{id}`: {error}\nSend another tag, or press Back.",
	RemoveFormHint:     "🗑 *Remove?* Set the options, then press Done.",
	NoContainers:       "There are no containers.",
	NoImages:           "There are no images.",
	DockerUnavailable:  "⚠️ Docker didn't answer: {error}",
	StreamStarting:     "⏳ Connecting to {name}…",
	Unauthorized: "⛔️ You're not allowed to use this bot.\n" +
		"Your chat ID is `{id}`. To allow it, start dockeroller with `-w {id}`.",

	ToastContainerStarted:  "Container started",
	ToastContainerStopped:  "Container stopped",
	ToastContainerPaused:   "Container paused",
	ToastContainerUnpaused: "Container unpaused",
	ToastContainerRemoved:  "Container removed",
	ToastImageRemoved:      "Image removed",
	ToastActionFailed:      "Failed: {error}",
	ToastListChanged:       "The list changed, showing it again",
	ToastInvalidButton:     "Invalid button 🤔",
	ToastOutdatedButton:    "This button is outdated, send /start",
	ToastFillFormPressDone: "Set the options and press Done",
	ToastNoContainers:      "There are no containers",
	ToastNoImages:          "There are no images",

	BtnContainers:    "Containers 📦",
	BtnImages:        "Images 🌄",
	BtnPrev:          "⬅ Prev",
	BtnNext:          "Next ➡",
	BtnBack:          "↩ Back",
	BtnStart:         "Start ▶️",
	BtnStop:          "Stop ⏹",
	BtnPause:         "Pause ⏸",
	BtnUnpause:       "Unpause ⏯",
	BtnRemove:        "Remove 🗑",
	BtnRename:        "Rename ✏️",
	BtnLogs:          "Logs 🪵",
	BtnStats:         "Stats 📊",
	BtnTag:           "Tag 🏷",
	BtnDone:          "Done ✅",
	BtnForce:         "Force: {value}",
	BtnRemoveVolumes: "Remove volumes: {value}",
	BtnPruneChildren: "Prune children: {value}",
}
