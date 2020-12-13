package routers

import (
	"oauth/internal/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter(u controllers.UserController) *gin.Engine {
	router := gin.Default()

	router.POST("/signup", u.AddUserController)
	router.POST("/login", u.LoginController)

	return router
}
