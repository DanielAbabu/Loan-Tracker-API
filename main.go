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

	logRepo := repositories.NewLogRepository(client)
	logUsecase := usecase.NewLogUsecase(logRepo)
	LogController := controllers.NewLogController(logUsecase)

	userRepo := repositories.NewUserRepository(client)
	userUsecase := usecase.NewUserUsecase(userRepo)
	UserController := controllers.NewUserController(userUsecase, logUsecase)

	loanRepo := repositories.NewLoanRepository(client)
	loanUsecase := usecase.NewLoanUsecase(loanRepo)
	LoanController := controllers.NewLoanController(loanUsecase, logUsecase)

	route := gin.Default()
	router.SetRouter(route, *UserController, *LoanController, *LogController, client)
	route.Run()
}
