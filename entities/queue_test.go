package entities_test

import (
	"testing"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {
	assert := assert.New(t)
	t.Run("pop priorities", func(t *testing.T) {
		// Arrange
		queue := entities.NewQueue()
		firstInput := "cont1"
		secondInput := "cont2"
		queue.Push(firstInput).Push(secondInput)

		// Act
		firstOutput, err1 := queue.Pop()
		secondOutput, err2 := queue.Pop()

		// Assert
		assert.Nil(err1)
		assert.Nil(err2)

		assert.Equal(firstInput, firstOutput)
		assert.Equal(secondInput, secondOutput)
	})

	t.Run("empty queue", func(t *testing.T) {
		// Arrange
		queue := entities.NewQueue()

		// Act
		poppedItem, err := queue.Pop()

		// Assert
		assert.NotNil(err)
		assert.Equal("", poppedItem)
	})

	t.Run("queue representation", func(t *testing.T) {
		queue := entities.NewQueue()
		queue.Push("cont1").Push("cont2")
		printableString := queue.String()

		assert.Equal("cont1\ncont2\n", printableString)
	})
}
