package usecase

import (
	"context"
	"errors"
	"loan-tracker/domain"
	"loan-tracker/infrastructure"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecases struct {
	UserRepo domain.UserRepository
}

func NewUserUsecase(Userrepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecases{
		UserRepo: Userrepo,
	}
}

func (uc *UserUsecases) RegisterUser(c context.Context, user *domain.User) error {
	return uc.UserRepo.RegisterUser(user)
}

func (uc *UserUsecases) VerifyUserEmail(c context.Context, token string) error {
	return uc.UserRepo.VerifyUserEmail(token)
}

func (uc *UserUsecases) LoginUser(c context.Context, user domain.User) (string, error) {
	return uc.UserRepo.LoginUser(user)
}

func (uc *UserUsecases) TokenRefresh(c context.Context, refreshToken string) (string, error) {
	token, err := infrastructure.TokenClaimer(refreshToken)
	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(*domain.JWTClaim)
	if !ok || claims.Exp < time.Now().Unix() {
		return "", errors.New("refresh token expired")
	}

	// Fetch user by ID
	var user domain.User
	user.ID, _ = primitive.ObjectIDFromHex(claims.UserID)
	ruser, err := uc.UserRepo.FindByID(user)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Verify the refresh token
	if ruser.RefreshToken != refreshToken {
		return "", errors.New("invalid refresh token, please login again")
	}

	// Generate new access token
	newAccessToken, err := infrastructure.TokenGenerator(user.ID, user.Email, true)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return newAccessToken, nil
}

func (uc *UserUsecases) UserProfile(c context.Context, user domain.User) (domain.ResponseUser, error) {
	return uc.UserRepo.UserProfile(user)
}

func (uc *UserUsecases) PasswordResetRequest(c context.Context, email string) error {
	return uc.UserRepo.PasswordResetRequest(email)
}

func (uc *UserUsecases) PasswordReset(c context.Context, token string, newPassword string) error {
	return uc.UserRepo.PasswordReset(token, newPassword)
}

func (uc *UserUsecases) GetAllUsers(c context.Context) ([]domain.ResponseUser, error) {
	return uc.UserRepo.GetAllUsers()
}

func (uc *UserUsecases) DeleteUser(c context.Context, user domain.User) error {
	return uc.UserRepo.DeleteUser(user)
}
