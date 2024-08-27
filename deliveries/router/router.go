package router

import (
	"loan-tracker/deliveries/controllers"
	"loan-tracker/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetRouter(router *gin.Engine, uc controllers.UserController, lc controllers.LoanController, loc controllers.LogController, client *mongo.Client) {
	router.POST("/users/register", uc.RegisterUser)
	router.GET("/users/verify-email", uc.VerifyUserEmail)
	router.POST("/users/login", uc.LoginUser)

	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(client))

	authRoutes.POST("/users/token/refresh", uc.TokenRefresh)
	authRoutes.POST("/users/password-update", uc.PasswordResetRequest)
	authRoutes.POST("/users/password-reset", uc.PasswordReset)
	authRoutes.GET("/users/profile", middleware.AuthMiddleware(client), uc.UserProfile)

	authRoutes.POST("/loans", lc.ApplyForLoan)
	authRoutes.GET("/loans/:id", lc.ViewLoanStatus)

	// Admin routes
	adminRoutes := authRoutes.Group("/admin")
	adminRoutes.Use(middleware.AdminMiddleware())

	adminRoutes.GET("/users", middleware.AuthMiddleware(client), uc.GetAllUsers)
	adminRoutes.DELETE("/users/:userid", middleware.AuthMiddleware(client), uc.DeleteUser)

	adminRoutes.GET("/loans", lc.ViewAllLoans)
	adminRoutes.PATCH("/loans/:id/status", lc.ApproveOrRejectLoan)
	adminRoutes.DELETE("/loans/:id", lc.DeleteLoan)
	adminRoutes.GET("/logs", loc.ViewSystemLogs)

}
