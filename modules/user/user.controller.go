package user

import (
	"v1/middleware"

	"github.com/gin-gonic/gin"
)

func UserController(r *gin.Engine) *gin.Engine {
	r.POST("/signup", Signup)
	r.POST("/login", Login)
	r.GET("/validate", middleware.Auth, ValidateToken)

	return r
}
