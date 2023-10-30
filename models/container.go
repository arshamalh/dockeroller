package models

import "fmt"

type Container struct {
	ID     string
	Name   string
	Image  string
	Status string
}

func (c Container) String() string {
	return fmt.Sprintf("%s - %s - %s - image: %s", c.ID, c.Name, c.Status, c.Image)
}
