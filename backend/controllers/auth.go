package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthController struct{}

func (a AuthController) Login(c *gin.Context) {
	var req loginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Welcome home!", "email": req.Email})
}

func (a AuthController) Test(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://google.com")
}
