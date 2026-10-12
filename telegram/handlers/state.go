package handlers

import (
	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/tbb"
)

// State is dockeroller's part of each chat's session. tbb persists it (in
// SQLite by default) whenever a handler changes it.
type State struct {
	// Containers and Images cache the last listing. Buttons carry an index
	// into them, so the list a user is paging through doesn't shift under
	// them while they look at it.
	Containers []*entities.Container `json:"containers,omitempty"`
	Images     []*entities.Image     `json:"images,omitempty"`

	// The container or image a rename/tag scene applies to.
	CurrentContainer *entities.Container `json:"current_container,omitempty"`
	CurrentImage     *entities.Image     `json:"current_image,omitempty"`

	// Options of the remove forms, toggled by their buttons before "Done".
	ContainerRemove entities.ContainerRemoveForm `json:"container_remove"`
	ImageRemove     entities.ImageRemoveForm     `json:"image_remove"`
}

// Ctx is the handler context.
type Ctx = tbb.Ctx[State]

// Scenes: the flows that expect a text answer.
const (
	SceneRenameContainer tbb.Scene = iota + 1
	SceneTagImage
)

// Questions inside those scenes.
const (
	QNewContainerName tbb.Question = 1
	QNewImageTag      tbb.Question = 1
)
