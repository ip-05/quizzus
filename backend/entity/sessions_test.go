package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	// Given
	wantGameId := uint(1)
	wantUserId := uint(1)
	wantInstId := uint(1)
	wantTime := time.Now()

	// When
	actual := NewSession(wantGameId, wantUserId, wantInstId)

	// Then
	assert.Equal(t, wantGameId, actual.GameId)
	assert.Equal(t, wantUserId, actual.UserId)
	assert.Equal(t, wantTime, actual.StartedAt)
}
