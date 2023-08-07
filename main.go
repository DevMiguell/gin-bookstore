package main

import (
	"net/http"
	config "v1/config"
	"v1/modules/book"
	"v1/modules/user"

	"github.com/gin-gonic/gin"
	// _ "github.com/99designs/gqlgen"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	book.BookController(r)
	user.UserController(r)

	http.ListenAndServe(":8080", r)
}
