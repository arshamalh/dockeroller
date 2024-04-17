package entities_test

import (
	"testing"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/stretchr/testify/assert"
)

func TestContainer(t *testing.T) {
	assert := assert.New(t)
	newContainer := entities.Container{
		ID:     "identified12",
		Name:   "Sweet Leonardo",
		Image:  "postgres:10.20-alpine3.15",
		Status: "Exited with 0",
		State:  entities.ContainerStateCreated,
	}

	assert.Equal(
		"identified12 - Sweet Leonardo - Exited with 0 - image: postgres:10.20-alpine3.15",
		newContainer.String(),
	)

	assert.Equal(false, newContainer.IsOn())
	newContainer.State = entities.ContainerStateRunning
	assert.Equal(true, newContainer.IsOn())
}
