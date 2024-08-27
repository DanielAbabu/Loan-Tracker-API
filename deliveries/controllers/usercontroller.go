package controllers

import (
	"loan-tracker/domain"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	Userusecase domain.UserUsecase
	LogUsecase  domain.LogUsecase
}

// Blog-controller constructor
func NewUserController(Usermgr domain.UserUsecase, logUsecase domain.LogUsecase) *UserController {
	return &UserController{
		Userusecase: Usermgr,
		LogUsecase:  logUsecase,
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

	// Log user registration
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "user_registration",
		Details:   "User registered with email: " + user.Email,
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging user registration:", logErr)
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

	// Log email verification
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "email_verification",
		Details:   "Email verified with token: " + token,
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging email verification:", logErr)
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

	// Log login attempt
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "login_attempt",
		Details:   "User login attempt with email: " + user.Email,
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging login attempt:", logErr)
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

	// Log token refresh
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "token_refresh",
		Details:   "Token refreshed with refresh token: " + tokenRequest.RefreshToken,
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging token refresh:", logErr)
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

	// Log user profile retrieval
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "user_profile_retrieval",
		Details:   "User profile retrieved for user ID: " + userID.Hex(),
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging user profile retrieval:", logErr)
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

	// Log password reset request
	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "password_reset_request",
		Details:   "Password reset requested for email: " + rest.Email,
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging password reset request:", logErr)
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

	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "password_reset_completion",
		Details:   "Password reset completed with token: " + token,
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging password reset completion:", logErr)
	}

	c.JSON(200, gin.H{"message": "password reset successful"})
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.Userusecase.GetAllUsers(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "get_all_users",
		Details:   "All users retrieved",
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging retrieval of all users:", logErr)
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

	logEntry := domain.Log{
		Timestamp: time.Now(),
		Type:      "user_deletion",
		Details:   "User deleted with ID: " + userID.Hex(),
	}
	if logErr := uc.LogUsecase.LogEvent(c, logEntry); logErr != nil {
		log.Println("Error logging user deletion:", logErr)
	}

	c.JSON(200, gin.H{"message": "user deleted successfully"})
}
