package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQuestion(t *testing.T) {
	body := CreateQuestion{
		Name: "What color is tomato?",
		Options: []CreateOption{
			{Name: "Red", Correct: true},
			{Name: "Green", Correct: false},
			{Name: "Blue", Correct: false},
			{Name: "Orange", Correct: false},
		},
	}

	actual, err := NewQuestion(body)
	assert.Nil(t, err)

	assert.Equal(t, body.Name, actual.Name)
	assert.Equal(t, body.Options[0].Name, actual.Options[0].Name)
	assert.Equal(t, body.Options[3].Name, actual.Options[3].Name)
	assert.Equal(t, body.Options[3].Correct, actual.Options[3].Correct)
}

func TestValidateQuestion(t *testing.T) {
	body := CreateQuestion{
		Name: "What color is tomato?",
		Options: []CreateOption{
			{Name: "Red", Correct: true},
			{Name: "Green", Correct: false},
			{Name: "Blue", Correct: false},
			{Name: "Orange", Correct: false},
		},
	}

	actual, err := NewQuestion(body)
	assert.Nil(t, err)

	actual.Options = []*Option{}
	errValidate := actual.Validate()

	assert.Contains(t, errValidate.Error(), "should be 4 options")
}
