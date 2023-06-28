package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	// Given
	wantLen := 9
	wantContain := "-"

	// When
	gotCode := GenerateCode()

	// Then
	assert.Equal(t, wantLen, len(gotCode))
	assert.Contains(t, gotCode, wantContain)
}

func TestGenerateToken(t *testing.T) {
	// Given
	wantContain := "."

	// When
	gotToken, err := GenerateToken(1, "test", "secret")

	// Then
	assert.Contains(t, gotToken, wantContain)
	assert.Nil(t, err)
}
