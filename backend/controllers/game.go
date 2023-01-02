package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameInfo struct {
	InviteCode string     `json:"code"`
	Topic      string     `json:"topic"`
	RoundTime  int        `json:"roundTime"`
	Points     float64    `json:"points"`
	Questions  []Question `json:"questions"`
}

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
	var gameInfo GameInfo
	var body CreateBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	gameInfo.InviteCode = generateCode()

	if len(body.Topic) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too long topic name"})
		return
	}
	gameInfo.Topic = body.Topic

	if body.RoundTime < 10 || body.RoundTime > 60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "round time should be over 10 or below 60 (seconds)"})
		return
	}
	gameInfo.RoundTime = body.RoundTime

	if body.Points <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "points should not be lower than 0"})
		return
	}
	gameInfo.Points = body.Points

	for _, v := range body.Questions {
		if len(v.Options) != 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "should be 4 options"})
			return
		}
		gameInfo.Questions = append(gameInfo.Questions, v)
	}

	c.JSON(http.StatusOK, gameInfo)
}

func (g GameController) FindByCode(c *gin.Context) {
	//var game GameInfo
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
