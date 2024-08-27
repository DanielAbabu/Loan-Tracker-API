package repositories

import (
	"context"
	"errors"
	"loan-tracker/domain"
	"loan-tracker/infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Col *mongo.Collection
}

func NewUserRepository(client *mongo.Client) domain.UserRepository {
	return &UserRepository{
		Col: client.Database("loan-tracker").Collection("users"),
	}
}

func (ur *UserRepository) RegisterUser(user *domain.User) error {
	filterUser := bson.M{"email": user.Email}
	var result domain.User
	err := ur.Col.FindOne(context.Background(), filterUser).Decode(&result)
	if err == nil {
		if result.IsVerified {
			return errors.New("user already exists")
		}
	}

	user.ID = primitive.NewObjectID()
	user.IsVerified = false
	user.IsAdmin = false

	password, err := infrastructure.PasswordHasher(user.Password)
	if err != nil {
		return errors.New(`password hashing failed`)
	}
	user.Password = password

	err = infrastructure.UserVerification(user.Email)
	if err != nil {
		return err
	}

	_, err = ur.Col.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil

}

func (ur *UserRepository) VerifyUserEmail(token string) error {
	email, err := infrastructure.VerifyToken(token)
	if err != nil {
		return err
	}

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"isverified": true}}

	_, err = ur.Col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) LoginUser(user domain.User) (string, error) {
	filter := bson.M{"email": user.Email}
	var fuser domain.User
	err := ur.Col.FindOne(context.TODO(), filter).Decode(&fuser)
	if err != nil {
		return "", err
	}

	if !fuser.IsVerified {
		return "", errors.New("user not verified")
	}

	check := infrastructure.PasswordComparator(fuser.Password, user.Password)
	if check != nil {
		return "", errors.New("invalid password")
	}

	accessToken, err := infrastructure.TokenGenerator(fuser.ID, fuser.Email, true)
	if err != nil {
		return "", errors.New("refreshtoken generation failed")
	}

	refreshToken, err := infrastructure.TokenGenerator(fuser.ID, fuser.Email, false)
	if err != nil {
		return "", errors.New("refreshtoken generation failed")
	}

	err = ur.TokenRefresh(fuser, refreshToken)
	if err != nil {
		return "", errors.New("refreshtoken update failed")
	}

	return accessToken, nil
}

func (ur *UserRepository) TokenRefresh(user domain.User, token string) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"refreshtoken": token,
		},
	}
	_, err := ur.Col.UpdateOne(context.TODO(), filter, update)
	return err
}

func (ur *UserRepository) UserProfile(user domain.User) (domain.ResponseUser, error) {
	filter := bson.M{"_id": user.ID}
	var result domain.ResponseUser
	err := ur.Col.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return domain.ResponseUser{}, err
	}
	return result, nil
}
func (ur *UserRepository) FindByID(user domain.User) (domain.User, error) {
	filter := bson.M{"_id": user.ID}
	var fuser domain.User
	err := ur.Col.FindOne(context.Background(), filter).Decode(&fuser)
	if err != nil {
		return domain.User{}, err
	}
	return fuser, nil
}
func (ur *UserRepository) PasswordResetRequest(email string) error {
	var user domain.User

	query := bson.M{"email": email}
	if err := ur.Col.FindOne(context.TODO(), query).Decode(&user); err != nil {
		return errors.New("user not found")
	}

	err := infrastructure.ForgotPasswordHandler(email)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) PasswordReset(token string, newPassword string) error {
	email, err := infrastructure.VerifyToken(token)
	if err != nil {
		return errors.New("token verification failed")
	}

	hashedPassword, err := infrastructure.PasswordHasher(newPassword)
	if err != nil {
		return err
	}

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": string(hashedPassword)}}

	_, err = ur.Col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetAllUsers() ([]domain.ResponseUser, error) {
	cur, err := ur.Col.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var users []domain.ResponseUser
	for cur.Next(context.Background()) {
		var user domain.ResponseUser
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) DeleteUser(user domain.User) error {
	filter := bson.M{"_id": user.ID}

	_, err := ur.Col.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
