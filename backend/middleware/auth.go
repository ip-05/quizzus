package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type AuthedUser struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing `token` cookie."})
			return
		}

		secret := []byte("sussyballs")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid `token` cookie."})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims["exp"].(float64))
			now := time.Now().Unix()
			if now > exp {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Expired `token` cookie."})
				return
			}

			authedUser := AuthedUser{
				Id:             claims["id"].(string),
				Name:           claims["name"].(string),
				Email:          claims["email"].(string),
				ProfilePicture: claims["profilePicture"].(string),
			}
			c.Set("authedUser", authedUser)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid `token` cookie."})
		}
	}
}
