package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ip-05/quizzus/config"
)

func WSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("token")
		if query == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing `token` parameter."})
			return
		}

		secret := []byte(cfg.Secrets.Jwt)
		token, err := jwt.Parse(query, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims["exp"].(float64))
			now := time.Now().Unix()
			if now > exp {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Expired token."})
				return
			}

			authedUser := AuthedUser{
				Id:   claims["id"].(uint),
				Name: claims["name"].(string),
			}
			c.Set("authedUser", authedUser)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token."})
		}
	}
}
