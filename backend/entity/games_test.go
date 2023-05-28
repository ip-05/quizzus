package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	createBody := CreateBody{
		Topic:     "My game",
		RoundTime: 10,
		Points:    3,
		Questions: []CreateQuestion{
			{
				Name: "What color is tomato?",
				Options: []CreateOption{
					{Name: "Red", Correct: true},
					{Name: "Green", Correct: false},
					{Name: "Blue", Correct: false},
					{Name: "Orange", Correct: false},
				},
			},
		},
	}

	t.Run("TestCreateOK", func(t *testing.T) {
		actual, err := NewGame(createBody, "123")
		assert.Nil(t, err)

		assert.Equal(t, createBody.Topic, actual.Topic)
		assert.Equal(t, createBody.RoundTime, actual.RoundTime)
		assert.Equal(t, createBody.Points, actual.Points)

		assert.Equal(t, createBody.Questions[0].Name, actual.Questions[0].Name)
		assert.Equal(t, createBody.Questions[0].Options[0].Name, actual.Questions[0].Options[0].Name)
		assert.Equal(t, createBody.Questions[0].Options[3].Name, actual.Questions[0].Options[3].Name)
	})

	t.Run("TestCreateErr", func(t *testing.T) {
		createBody.Points = 0
		actual, err := NewGame(createBody, "123")
		assert.Nil(t, actual)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "points should not be lower than 0")
	})
}

func TestValidateGame(t *testing.T) {
	createBody := CreateBody{
		Topic:     "My game",
		RoundTime: 10,
		Points:    3,
		Questions: []CreateQuestion{
			{
				Name: "What color is tomato?",
				Options: []CreateOption{
					{Name: "Red", Correct: true},
					{Name: "Green", Correct: false},
					{Name: "Blue", Correct: false},
					{Name: "Orange", Correct: false},
				},
			},
		},
	}

	actual, err := NewGame(createBody, "123")
	assert.Nil(t, err)

	t.Run("TestTopic", func(t *testing.T) {
		actual.Topic = strings.Repeat(".", 129)
		errValidate := actual.Validate()
		assert.Contains(t, errValidate.Error(), "too long topic name")
		actual.Topic = "Topic"
	})

	t.Run("TestRoundTime", func(t *testing.T) {
		actual.RoundTime = 0
		errValidate := actual.Validate()
		assert.Contains(t, errValidate.Error(), "round time should be over 10 or below 60 (seconds)")
		actual.RoundTime = 10
	})

	t.Run("TestPoints", func(t *testing.T) {
		actual.Points = 0
		errValidate := actual.Validate()
		assert.Contains(t, errValidate.Error(), "points should not be lower than 0")
		actual.Points = 10
	})
	t.Run("TestQuestions", func(t *testing.T) {
		actual.Questions = []*Question{}
		errValidate := actual.Validate()
		assert.Contains(t, errValidate.Error(), "should be at least 1 question")
	})
}

func TestGenerateCode(t *testing.T) {
	// Given
	wantLen := 9
	wantContain := "-"

	// When
	gotCode := generateCode()

	// Then
	assert.Equal(t, wantLen, len(gotCode))
	assert.Contains(t, gotCode, wantContain)
}
