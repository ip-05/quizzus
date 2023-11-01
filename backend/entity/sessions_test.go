package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	// Given
	wantGameID := uint(1)
	wantUserID := uint(1)
	wantInstID := uint(1)

	// When
	actual := NewSession(wantGameID, wantUserID, wantInstID)

	// Then
	assert.Equal(t, wantGameID, actual.GameID)
	assert.Equal(t, wantUserID, actual.UserID)
}
