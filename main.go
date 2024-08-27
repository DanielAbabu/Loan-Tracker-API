package main

import (
	"loan-tracker/deliveries/controllers"
	"loan-tracker/deliveries/router"
	"loan-tracker/infrastructure"
	"loan-tracker/repositories"
	"loan-tracker/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	client := infrastructure.MongoDBInit()
	userRepo := repositories.NewUserRepository(client)
	userUsecase := usecase.NewUserUsecase(userRepo)
	UserController := controllers.NewUserController(userUsecase)

	route := gin.Default()
	router.SetRouter(route, *UserController, client)
	route.Run()
}
