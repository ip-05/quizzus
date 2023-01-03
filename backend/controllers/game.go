package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/models"
)

type Question struct {
	Name    string   `json:"name"`
	Options []Option `json:"options"`
}

type Option struct {
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}

type CreateBody struct {
	Topic     string     `json:"topic"`
	RoundTime int        `json:"roundTime"`
	Points    float64    `json:"points"`
	Questions []Question `json:"questions"`
}

type GameController struct{}

func (g GameController) CreateGame(c *gin.Context) {
	var game models.Game

	var body CreateBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	game.InviteCode = generateCode()

	if len(body.Topic) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too long topic name"})
		return
	}
	game.Topic = body.Topic

	if body.RoundTime < 10 || body.RoundTime > 60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "round time should be over 10 or below 60 (seconds)"})
		return
	}
	game.RoundTime = body.RoundTime

	if body.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "points should not be lower than 0"})
		return
	}
	game.Points = body.Points

	//models.DB.Create(&game)

	for _, v := range body.Questions {
		if len(v.Options) != 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "should be 4 options"})
			return
		}

		question := models.Question{Name: v.Name}
		//models.DB.Create(&question)
		game.Questions = append(game.Questions, question)

		for _, j := range v.Options {
			option := models.Option{Name: j.Name, Correct: j.Correct}
			//models.DB.Create(&option)
			question.Options = append(question.Options, option)
		}
	}
	models.DB.Create(&game.Questions)
	models.DB.Create(&game)

	models.DB.First(&game, game.Id).Preload("Questions")
	c.JSON(http.StatusOK, game)
}

func (g GameController) FindByCode(c *gin.Context) {
	//var game game
	c.JSON(http.StatusBadRequest, gin.H{"message": "not yet implemented"})
}

func generateCode() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	code := hex.EncodeToString(bytes)

	return fmt.Sprintf("%s-%s", code[:4], code[4:])
}
