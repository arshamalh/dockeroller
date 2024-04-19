package entities

import "fmt"

type ContainerState string

const (
	ContainerStateCreated    ContainerState = "created"
	ContainerStateRunning    ContainerState = "running"
	ContainerStatePaused     ContainerState = "paused"
	ContainerStateRestarting ContainerState = "restarting"
	ContainerStateRemoving   ContainerState = "removing"
	ContainerStateExited     ContainerState = "exited"
	ContainerStateDead       ContainerState = "dead"
)

type Container struct {
	ID         string
	Name       string
	Image      string
	Status     string
	State      ContainerState
	RemoveForm *ContainerRemoveForm
}

func (c Container) String() string {
	return fmt.Sprintf("%s - %s - %s - image: %s", c.ID, c.Name, c.Status, c.Image)
}

func (c Container) IsOn() bool {
	return c.State == ContainerStateRunning
}
