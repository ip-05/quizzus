package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOption(t *testing.T) {
	body := CreateOption{Name: "Red", Correct: true}

	actual, err := NewOption(body)
	assert.Nil(t, err)
	assert.Equal(t, body.Name, actual.Name)
	assert.Equal(t, body.Correct, actual.Correct)
}
