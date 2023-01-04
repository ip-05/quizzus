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
	Id      uint     `json:"id"`
	GameID  uint     `json:"gameId"`
}

type Option struct {
	Name       string `json:"name"`
	Correct    bool   `json:"correct"`
	Id         uint   `json:"id"`
	QuestionID uint   `json:"questionId"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too long topic name"})
		return
	}
	game.Topic = body.Topic

	if body.RoundTime < 10 || body.RoundTime > 60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Round time should be over 10 or below 60 (seconds)"})
		return
	}
	game.RoundTime = body.RoundTime

	if body.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Points should not be lower than 0"})
		return
	}
	game.Points = body.Points

	models.DB.Create(&game)

	for _, v := range body.Questions {
		if len(v.Options) != 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Should be 4 options"})
			return
		}
		question := models.Question{Name: v.Name, GameID: game.Id}
		models.DB.Create(&question)

		for _, j := range v.Options {
			option := models.Option{Name: j.Name, Correct: j.Correct, QuestionID: question.Id}
			models.DB.Create(&option)
		}
	}
	c.JSON(http.StatusOK, game)
}

func (g GameController) GetById(c *gin.Context) {
	var game models.Game
	id := c.Param("id")
	models.DB.Preload("Questions.Options").First(&game, id)
	c.JSON(http.StatusOK, game)
}

func (g GameController) GetByCode(c *gin.Context) {
	var game models.Game

	code := c.Param("code")
	models.DB.Preload("Questions.Options").Where("invite_code = ?", code).First(&game)

	if game.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	c.JSON(http.StatusOK, game)
}

// func (g GameController) Update(c *gin.Context) {
// 	var game models.Game
// 	var body CreateBody

// 	if err := c.BindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err})
// 		return
// 	}

// 	id := c.Param("id")
// 	models.DB.Preload("Questions.Options").Where("id = ?", id).First(&game)

// 	// if err := models.DB.Preload("Questions.Options").Where("id = ?", id).First(&game).Error; err != nil {
// 	// 	c.AbortWithStatus(404)
// 	// 	fmt.Println(err)
// 	// }
// 	// c.BindJSON(&game)
// 	// models.DB.Save(&game)
// 	// c.JSON(200, game)
// 	models.DB.Model(&game).Where("id = ?", id).Updates(models.Game{Topic: body.Topic, RoundTime: body.RoundTime, Points: body.Points})

// 	for _, q := range body.Questions {
// 		var question models.Question
// 		qId := q.Id
// 		models.DB.Model(&question).Where("id = ?", qId).Updates(models.Question{Name: q.Name})

// 		for _, o := range q.Options {
// 			var option models.Option
// 			oId := o.Id
// 			models.DB.Model(&option).Where("id = ?", oId).Updates(models.Option{Name: o.Name, Correct: o.Correct})

// 		}
// 		//v.Name = body.Questions
// 	}
// 	//models.DB.Save(&game)
// 	c.JSON(http.StatusOK, game)
// }

// func (g GameController) Delete(c *gin.Context) {
// 	//var game models.Game
// }

func generateCode() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	code := hex.EncodeToString(bytes)

	return fmt.Sprintf("%s-%s", code[:4], code[4:])
}
