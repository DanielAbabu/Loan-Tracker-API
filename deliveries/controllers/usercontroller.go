package controllers

import (
	"loan-tracker/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	Userusecase domain.UserUsecase
}

// Blog-controller constructor
func NewUserController(Usermgr domain.UserUsecase) *UserController {
	return &UserController{
		Userusecase: Usermgr,
	}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" {
		c.JSON(400, gin.H{"error": "Please provide an email"})
		return
	}

	if user.Password == "" {
		c.JSON(400, gin.H{"error": "Please provide a password"})
		return
	}

	if user.UserName == "" {
		c.JSON(400, gin.H{"error": "Please provide a username"})
		return
	}

	user.IsAdmin = false

	err := uc.Userusecase.RegisterUser(c, &user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User registered successfully"})
}

func (uc *UserController) VerifyUserEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"error": "Please provide a token"})
		return
	}

	err := uc.Userusecase.VerifyUserEmail(c, token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Email verified successfully"})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" {
		c.JSON(400, gin.H{"error": "Please provide an email"})
		return
	}

	if user.Password == "" {
		c.JSON(400, gin.H{"error": "Please provide a password"})
		return
	}

	token, err := uc.Userusecase.LoginUser(c, user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user successfully logged in", "token": token})
}

func (uc *UserController) TokenRefresh(c *gin.Context) {
	var tokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	newToken, err := uc.Userusecase.TokenRefresh(c, tokenRequest.RefreshToken)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "token refreshed successfully", "token": newToken})
}

func (uc *UserController) UserProfile(c *gin.Context) {

	// Get the user id from the context
	userid := c.GetString("userid")
	if userid == "" {
		c.JSON(400, gin.H{"error": "Please provide a userid"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userid)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user domain.User
	user.ID = userID
	fuser, err := uc.Userusecase.UserProfile(c, user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user profile retrieved successfully", "user": fuser})
}

func (uc *UserController) PasswordResetRequest(c *gin.Context) {
	var rest domain.RestRequest
	if err := c.ShouldBindJSON(&rest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if rest.Email == "" {
		c.JSON(400, gin.H{"error": "Please provide an email"})
		return
	}

	err := uc.Userusecase.PasswordResetRequest(c, rest.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "password reset request successful"})
}

func (uc *UserController) PasswordReset(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"error": "Please provide a token"})
		return
	}

	var rest domain.RestRequest
	if err := c.ShouldBindJSON(&rest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if rest.NewPassword == "" {
		c.JSON(400, gin.H{"error": "Please provide a new password"})
		return
	}

	err := uc.Userusecase.PasswordReset(c, token, rest.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "password reset successful"})
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.Userusecase.GetAllUsers(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "users retrieved successfully", "users": users})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userid := c.Param("userid")
	if userid == "" {
		c.JSON(400, gin.H{"error": "Please provide a userid"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userid)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user domain.User
	user.ID = userID
	err = uc.Userusecase.DeleteUser(c, user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user deleted successfully"})
}
