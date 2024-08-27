package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserName     string             `json:"username,omitempty"`
	Email        string             `json:"email,omitempty"`
	Password     string             `json:"password,omitempty"`
	IsAdmin      bool               `json:"isadmin,omitempty"`
	RefreshToken string             `json:"refreshtoken,omitempty"`
	IsVerified   bool               `bson:"isverified,omitempty" json:"isverified,omitempty"`
}

type ResponseUser struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserName   string             `json:"username,omitempty"`
	Email      string             `json:"email,omitempty"`
	IsAdmin    bool               `bson:"isadmin,omitempty" json:"isadmin"`
	IsVerified bool               `bson:"isverified,omitempty" json:"isverified"`
}

type RestRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"password"`
}

type UserUsecase interface {
	RegisterUser(c context.Context, user *User) error
	VerifyUserEmail(c context.Context, token string) error
	LoginUser(c context.Context, user User) (string, error)
	TokenRefresh(c context.Context, token string) (string, error)
	UserProfile(c context.Context, user User) (ResponseUser, error)
	PasswordResetRequest(c context.Context, email string) error
	PasswordReset(c context.Context, token string, newPassword string) error
	GetAllUsers(c context.Context) ([]ResponseUser, error)
	DeleteUser(c context.Context, user User) error
}

type UserRepository interface {
	RegisterUser(user *User) error
	VerifyUserEmail(token string) error
	LoginUser(user User) (string, error)
	TokenRefresh(user User, token string) error
	UserProfile(user User) (ResponseUser, error)
	PasswordResetRequest(email string) error
	PasswordReset(token string, newPassword string) error
	GetAllUsers() ([]ResponseUser, error)
	DeleteUser(user User) error
	FindByID(user User) (User, error)
}
