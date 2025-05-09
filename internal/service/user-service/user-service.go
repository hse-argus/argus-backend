package userservice

import (
	"argus-backend/internal/logger"
	"argus-backend/internal/repository/user"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

type UserServiceInterface interface {
	Login(login string, password string) (string, error)
	Register(creds user.User) (string, error)
}

type UserService struct {
	userRepo *user.UserRepository
}

func NewUserService(repo *user.UserRepository) UserServiceInterface {
	return &UserService{
		userRepo: repo,
	}
}

func (us *UserService) Login(login string, password string) (string, error) {
	existedUser, err := us.userRepo.GetUser(login, password)
	if err != nil {
		return "", err
	}

	return createJWTToken(existedUser)
}

func (us *UserService) Register(creds user.User) (string, error) {
	id, err := us.userRepo.AddUser(creds)
	if err != nil {
		return "", err
	}
	creds.Id = id

	return createJWTToken(&creds)
}

func createJWTToken(user *user.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":  user.Name,
		"login": user.Login,
		"id":    user.Id,
	})
	token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		logger.Error(fmt.Sprintf("error creating jtw token: %v", err))
		return "", err
	}
	return token, nil
}
