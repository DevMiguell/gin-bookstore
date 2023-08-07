package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"v1/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func generateMockToken(userID uint) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))

	return tokenString
}

func setupTestRouter(mockUserID uint) *gin.Engine {
	config.ConnectDatabase()
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if mockUserID != 0 {
			mockToken := generateMockToken(mockUserID)
			c.Request.Header.Set("Authorization", "Bearer "+mockToken)
		}
		c.Next()
	})

	r.POST("/signup", Signup)
	r.POST("/login", Login)
	r.GET("/validate", ValidateToken)

	return r
}

func TestSignup(t *testing.T) {
	config.ConnectDatabase()
	defer config.DB.Close()

	router := setupTestRouter(0)

	// Prepare request data
	input := CreateUserInput{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password",
	}

	inputJSON, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, input.Email, response["data"]["email"])
	assert.Equal(t, input.Username, response["data"]["username"])
}

func TestLogin(t *testing.T) {
	config.ConnectDatabase()
	defer config.DB.Close()

	router := setupTestRouter(0)
	input := LoginUserInput{
		Email:    "test@example.com",
		Password: "password",
	}

	inputJSON, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestValidateToken(t *testing.T) {
	config.ConnectDatabase()
	defer config.DB.Close()

	mockUserID := uint(123)
	router := setupTestRouter(mockUserID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/validate", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
