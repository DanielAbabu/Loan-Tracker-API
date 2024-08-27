package router

import (
	"loan-tracker/deliveries/controllers"
	"loan-tracker/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetRouter(router *gin.Engine, uc controllers.UserController, client *mongo.Client) {
	router.POST("/users/register", uc.RegisterUser)
	router.GET("/users/verify-email", uc.VerifyUserEmail)
	router.POST("/users/login", uc.LoginUser)
	router.POST("/users/token/refresh", uc.TokenRefresh)
	router.POST("/users/password-update", uc.PasswordResetRequest)
	router.POST("/users/password-reset", uc.PasswordReset)
	router.GET("/users/profile", middleware.AuthMiddleware(client), uc.UserProfile)

	router.GET("/admin/users", middleware.AuthMiddleware(client), uc.GetAllUsers)
	router.DELETE("/admin/users/:userid", middleware.AuthMiddleware(client), uc.DeleteUser)
}
