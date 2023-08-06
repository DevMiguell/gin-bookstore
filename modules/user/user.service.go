package user

import (
	"net/http"
	"os"
	"time"
	models "v1/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserWithoutPassword struct {
	ID       int    `json:"ID"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func Signup(c *gin.Context) {
	var input CreateUserInput

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	user := models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: string(hash),
	}

	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"email":    user.Email,
		"username": user.Username,
	}})
}

func Login(c *gin.Context) {
	var input LoginUserInput

	if c.Bind(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user models.User
	models.DB.First(&user, "email = ?", input.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password or email",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // TODO: optimize this for 30 days if necessary
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to create a token",
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func ValidateToken(c *gin.Context) {
	user, _ := c.Get("user")

	userWithPassword, ok := user.(models.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user data",
		})
		return
	}

	userWithoutPassword := UserWithoutPassword{
		ID:       int(userWithPassword.ID),
		Email:    userWithPassword.Email,
		Username: userWithPassword.Username,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":          userWithoutPassword,
		"authenticated": true,
	})
}
