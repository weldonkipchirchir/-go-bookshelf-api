package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/weldonkipchirchir/go/bookshelf-api/src/config"
	"github.com/weldonkipchirchir/go/bookshelf-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("username already taken")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func RegisterUser(user *models.User) error {
	// check if username is already taken

	existingUser, err := getUserByEmail(user.Email)
	if existingUser != nil {
		return ErrEmailTaken
	}

	if err != nil {
		return err
	}

	//hash password
	password := HashPassword(user.Password)
	user.Password = password

	user.Password = string(password)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	_, err = config.DB.Collection("users").InsertOne(ctx, user)
	return err
}

func LoginUser(email, password string) (*models.User, error) {
	user, err := getUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func getUserByEmail(email string) (*models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
