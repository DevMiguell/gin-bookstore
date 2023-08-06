package main

import (
	"net/http"
	models "v1/config"
	"v1/modules/book"
	"v1/modules/user"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()

	r := gin.Default()

	book.BookController(r)
	user.UserController(r)

	http.ListenAndServe(":8080", r)
}
