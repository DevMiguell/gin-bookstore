package book

import (
	"github.com/gin-gonic/gin"
)

func BookController(r *gin.Engine) *gin.Engine {
	r.GET("/books", FindBooks)
	r.GET("/book/:id", FindBook)
	r.POST("/book", CreateBook)
	r.PATCH("/book/:id", UpdateBook)
	r.DELETE("/book/:id", DeleteBook)

	return r
}
