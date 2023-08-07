package book

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// "time"
	"v1/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func setupTest() {
	config.ConnectDatabase()
	config.DB.Begin()
}

func teardownTest() {
	config.DB.Rollback()
	config.DB.Close()
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

	r.GET("/books", FindBooks)
	r.GET("/books/:id", FindBook)
	r.POST("/books", CreateBook)
	r.PUT("/books/:id", UpdateBook)
	r.DELETE("/books/:id", DeleteBook)
	return r
}

func MockFindBooks() {
	mockBooks := []config.Book{
		{Title: "Book 1"},
		{Title: "Book 2"},
		{Title: "Book 3"},
	}

	for _, book := range mockBooks {
		copyBook := book // Criar uma cópia do livro
		config.DB.Create(&copyBook)
	}
}

func generateMockToken(userID uint) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))

	return tokenString
}

func TestFindBooks(t *testing.T) {
	setupTest()
	defer teardownTest()

	MockFindBooks()

	mockUserID := uint(123)
	router := setupTestRouter(mockUserID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]config.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedBooks := []config.Book{
		{Title: "Book 1"},
		{Title: "Book 2"},
		{Title: "Book 3"},
	}

	assert.Equal(t, expectedBooks[0].Title, response["data"][0].Title)
}

func TestFindBook(t *testing.T) {
	setupTest()
	defer teardownTest()

	MockFindBooks()

	mockUserID := uint(123)
	router := setupTestRouter(mockUserID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]config.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedBook := config.Book{Title: "Book 1"} // Compare apenas o campo Title
	assert.Equal(t, expectedBook.Title, response["data"].Title)
}

func TestCreateBook(t *testing.T) {
	setupTest()
	defer teardownTest()

	mockUserID := uint(123)
	router := setupTestRouter(mockUserID)

	// Prepare request data
	input := CreateBookInput{
		Title:  "New Book",
		Author: "Author Name",
	}

	inputJSON, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]config.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Compare individual fields of the expected and actual books
	assert.Equal(t, input.Title, response["data"].Title)
	assert.Equal(t, input.Author, response["data"].Author)
}

func TestUpdateBook(t *testing.T) {
	setupTest()
	defer teardownTest()

	MockFindBooks()

	mockUserID := uint(123)
	router := setupTestRouter(mockUserID)

	input := UpdateBookInput{
		Title: "Updated Book Title",
	}

	inputJSON, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/1", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]config.Book
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Compare individual fields of the expected and actual books
	assert.Equal(t, input.Title, response["data"].Title)
}

func TestDeleteBook(t *testing.T) {
	setupTest()
	defer teardownTest()

	MockFindBooks()

	mockUserID := uint(123)
	router := setupTestRouter(mockUserID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/1", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]bool
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["data"])
}
