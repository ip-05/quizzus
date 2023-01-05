package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateQuestion struct {
	Name    string         `json:"name"`
	Options []CreateOption `json:"options"`
}

type CreateOption struct {
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}

type CreateBody struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"roundTime"`
	Points    float64          `json:"points"`
	Questions []CreateQuestion `json:"questions"`
}

type UpdateQuestion struct {
	Id      uint           `json:"id"`
	Name    string         `json:"name"`
	Options []UpdateOption `json:"options"`
}

type UpdateOption struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}

type UpdateBody struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"roundTime"`
	Points    float64          `json:"points"`
	Questions []UpdateQuestion `json:"questions"`
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

	for _, v := range body.Questions {
		if len(v.Options) != 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Should be 4 options"})
			return
		}
		question := models.Question{Name: v.Name, GameID: game.Id}

		for _, j := range v.Options {
			option := models.Option{Name: j.Name, Correct: j.Correct, QuestionID: question.Id}
			question.Options = append(question.Options, option)
		}
		game.Questions = append(game.Questions, question)
	}

	models.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&game)

	c.JSON(http.StatusOK, game)
}

func (g GameController) Get(c *gin.Context) {
	var game models.Game

	id, _ := strconv.Atoi(c.Query("id"))
	code := c.Query("invite_code")
	models.DB.Preload("Questions.Options").Where("invite_code = ? or id = ?", code, id).First(&game)

	if game.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	c.JSON(http.StatusOK, game)
}

func (g GameController) Update(c *gin.Context) {
	var game models.Game
	var body UpdateBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	id, _ := strconv.Atoi(c.Query("id"))
	code := c.Query("invite_code")
	models.DB.Preload("Questions.Options").Where("invite_code = ? or id = ?", code, id).First(&game)

	if game.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	game.Topic = body.Topic
	game.RoundTime = body.RoundTime
	game.Points = body.Points

	ids := make(map[uint]int)
	for _, y := range game.Questions {
		ids[y.Id] += 1
	}

	for i, x := range body.Questions {
		ids[x.Id] += 1

		fmt.Println(ids[x.Id])
		if ids[x.Id] == 2 {
			game.Questions[i].Name = x.Name

			for j := 0; j < 4; j++ {
				game.Questions[i].Options[j].Name = x.Options[j].Name
				game.Questions[i].Options[j].Correct = x.Options[j].Correct
			}
		} else {
			question := models.Question{

				Name: x.Name,
			}

			for i := 0; i < 4; i++ {
				question.Options = append(question.Options, models.Option{Name: x.Options[i].Name, Correct: x.Options[i].Correct})
			}

			game.Questions = append(game.Questions, question)
		}
	}

	for _, v := range ids {
		if v == 1 {
			// for j, v2 := range game.Questions {
			// 	if v2.Id == i {
			// 		game.Questions = append(game.Questions[:j], game.Questions[j+1:]...)
			// 	}
			// }
			// models.DB.Unscoped().Where("id = ?", i).Delete(&models.Question{})
			// DELETE HERE
			//models.DB.Select("Options").Unscoped().Delete(&models.Question{}, i)
		}
	}

	models.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&game)

	c.JSON(http.StatusOK, game)
}

func (g GameController) Delete(c *gin.Context) {
	var game models.Game

	id, _ := strconv.Atoi(c.Query("id"))
	code := c.Query("invite_code")
	models.DB.Preload("Questions.Options").Where("invite_code = ? or id = ?", code, id).First(&game)

	if game.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	for _, v := range game.Questions {
		models.DB.Select(clause.Associations).Unscoped().Delete(&v)
	}
	models.DB.Select(clause.Associations).Unscoped().Delete(&game)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}

func generateCode() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	code := hex.EncodeToString(bytes)

	return fmt.Sprintf("%s-%s", code[:4], code[4:])
}
