package session_test

import (
	"github.com/arshamalh/dockeroller/session"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSession(t *testing.T) {
	assert := assert.New(t)
	t.Run("new session", func(t *testing.T) {
		s := session.New()

		var userID int64 = 2
		userData := s.Get(userID)

		assert.NotNil(userData)
	})
}
