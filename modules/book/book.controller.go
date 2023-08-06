package book

import (
	"v1/middleware"

	"github.com/gin-gonic/gin"
)

func BookController(r *gin.Engine) *gin.Engine {
	r.GET("/books", FindBooks)
	r.GET("/book/:id", FindBook)
	r.POST("/book", middleware.Auth, CreateBook)
	r.PATCH("/book/:id", middleware.Auth, UpdateBook)
	r.DELETE("/book/:id", middleware.Auth, DeleteBook)

	return r
}
