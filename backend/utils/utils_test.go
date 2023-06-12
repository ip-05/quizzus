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
