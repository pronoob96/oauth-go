package controllers

import (
	"log"
	"net/http"
	"oauth/config"
	"oauth/internal/services"
	"oauth/pkg/dto"
	"oauth/utils"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	AddUserController(*gin.Context)
	LoginController(ctx *gin.Context)
}

type userController struct {
	conf     *config.YamlConfig
	services services.UserService
}

func NewUserController(services services.UserService, conf *config.YamlConfig) UserController {
	return &userController{
		conf:     conf,
		services: services,
	}
}

func (c *userController) AddUserController(ctx *gin.Context) {
	var userRequest dto.User

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	log.Println("User Request in Controller: ", userRequest)
	userResponse, err := c.services.AddUserService(ctx, &userRequest)
	if err != nil {
		log.Println(err)
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Response in Controller: ", userResponse)

	utils.RespondJSON(ctx.Writer, http.StatusOK, userResponse)
}

func (c *userController) LoginController(ctx *gin.Context) {
	var userRequest dto.User

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	log.Println("User Request in Controller: ", userRequest)
	token, err := c.services.LoginService(ctx, &userRequest)
	if err != nil {
		log.Println(err)
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(ctx.Writer, http.StatusOK, token)
}
