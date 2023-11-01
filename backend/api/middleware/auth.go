package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ip-05/quizzus/config"
)

type AuthedUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing `Authorization` header."})
			return
		}

		tokenString := strings.Split(header, " ")[1]

		secret := []byte(cfg.Secrets.Jwt)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims["exp"].(float64))
			now := time.Now().Unix()
			if now > exp {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Expired token."})
				return
			}

			authedUser := AuthedUser{
				ID:   uint(claims["id"].(float64)),
				Name: claims["name"].(string),
			}
			ctx.Set("authedUser", authedUser)
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token."})
		}
	}
}
