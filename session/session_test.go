package session_test

import (
	"testing"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/session"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	assert := assert.New(t)
	t.Run("new session", func(t *testing.T) {
		s := session.New()

		var userID int64 = 2
		userData := s.Get(userID)

		assert.NotNil(userData)
	})

	t.Run("scene", func(t *testing.T) {
		s := session.New()
		given := entities.SceneRenameContainer

		var userID int64 = 2
		userData := s.Get(userID)
		assert.NotNil(userData)

		userData.SetScene(given)
		got := userData.GetScene()

		assert.Equal(given, got)
	})
}
