package service

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepository repository.UserRepository
}

// func (serv UserService) SignUp(input any) {} // any and interface don't enforce types, beneficial if it's a changing input
// func (serv UserService) SignUp(input interface) {}
func (userService UserService) Register(userInfo dto.RegisterDTO) (string, error) {
	log.Println(userInfo)

	newUser := domain.User{
		Uuid:      uuid.New(),
		Email:     userInfo.Email,
		Password:  userInfo.Password,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Phone:     userInfo.Phone,
	}
	err := userService.UserRepository.CreateUser(&newUser)
	if err != nil {
		return "", err
	}

	fakeUserToken := fmt.Sprintf("%v, %v, %v", newUser.Uuid, newUser.Email, newUser.UserType)
	return fakeUserToken, err
}

func (userService UserService) Login(attempt dto.LoginDTO) (string, error) {
	foundUser, err := userService.findUserByEmail(attempt.Email)
	if err != nil {
		return "", errors.New("user not found: " + err.Error())
	}

	if foundUser.Password != attempt.Password {
		return "", errors.New("wrong password")
	}

	foundUser.LastLogin = time.Now()
	err = userService.UserRepository.UpdateUser(foundUser)

	fakeUserToken := fmt.Sprintf("logged in as %v, %v, %v", foundUser.Uuid, foundUser.Email, foundUser.UserType)
	return fakeUserToken, nil
}

func (userService UserService) findUserByEmail(email string) (*domain.User, error) {
	foundUser, err := userService.UserRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (userService UserService) GetVerificationCode(attempt domain.User) (int, error) {
	return 0, nil
}

func (userService UserService) VerifyCode(id uint, code int) error {
	return nil
}

func (userService UserService) CreateProfile(id uint, input any) error {
	return nil
}

func (userService UserService) GetProfile(id uint) (*domain.User, error) {
	return nil, nil
}

func (userService UserService) UpdateProfile(id uint, input any) error {
	return nil
}

func (userService UserService) BecomeSeller(id uint, input any) error {
	return nil
}

func (userService UserService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}

func (userService UserService) CreateCart(id uint, input any) error {
	return nil
}

func (userService UserService) CreateOrder(id uint, input any) error {
	return nil
}

func (userService UserService) GetOrders(id uint, input any) error {
	return nil
}

func (userService UserService) GetOrderById(id uint, input any) error {
	return nil
}
