package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"
	config "v1/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	if tokenString == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{
			"message": "Invalid or expired token 2",
		})
		c.Abort()
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	id := int(claims["sub"].(float64))

	if time.Now().After(expirationTime) {
		c.JSON(401, gin.H{
			"message": "Token has expired",
		})
		c.Abort()
		return
	}

	var user config.User
	config.DB.First(&user, "id = ?", id)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User is not authorized",
		})
		return
	}

	c.Set("user", user)

	c.Next()
}
