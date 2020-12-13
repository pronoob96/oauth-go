package main

import (
	"log"
	"oauth/config"
	"oauth/internal/controllers"
	"oauth/internal/repository"
	"oauth/internal/routers"
	"oauth/internal/services"
	"oauth/redis"
	"oauth/utils"
	"strconv"
)

func main() {

	configYaml := config.GetYamlServiceConfig()
	if configYaml == nil {
		log.Println("Error getting configuration!")
		return
	}
	utils.StartUp(&configYaml.Datasource)
	redis.InitializeRedis(&configYaml.Redis)

	port := strconv.FormatInt(configYaml.Server.Port, 10)

	log.Println(port)

	_, db := utils.GetDBConnection()
	usersCollection := db.Collection("users")

	userRepo := repository.NewUserRepo(usersCollection)

	userServices := services.NewUsersService(userRepo)

	userControllers := controllers.NewUserController(userServices, configYaml)

	router := routers.NewRouter(userControllers)

	router.Run(":" + port)
}
