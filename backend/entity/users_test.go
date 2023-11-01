package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	// Given
	body := &CreateUser{
		GoogleID:   "123",
		DiscordID:  "456",
		TelegramID: "789",
		Picture:    "picture.com",
		Name:       "test",
	}

	t.Run("TestOK", func(t *testing.T) {
		// When
		actual, err := NewUser(body)

		// Then
		assert.Equal(t, body.DiscordID, actual.DiscordID)
		assert.Nil(t, err)
	})

	t.Run("TestInvalidName", func(t *testing.T) {
		body.Name = "t"

		// When
		actual, err := NewUser(body)

		// Then
		assert.Nil(t, actual)
		assert.NotNil(t, err)
	})
}
