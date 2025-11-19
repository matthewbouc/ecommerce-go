package service

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"errors"
	"fmt"
	"log"
	"time"
)

type UserService struct {
	UserRepository repository.UserRepository
}

// func (serv UserService) SignUp(input any) {} // any and interface don't enforce types, beneficial if it's a changing input
// func (serv UserService) SignUp(input interface) {}
func (userService UserService) Register(userInfo dto.RegisterDTO) (string, error) {
	log.Println(userInfo)

	newUser := domain.User{
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

	holderToken := fmt.Sprintf("%v, %v, %v", newUser.Uuid, newUser.Email, newUser.UserType)
	return holderToken, err
}

func (userService UserService) Login(attempt dto.LoginDTO) (string, error) {
	foundUser, err := userService.findUserByEmail(attempt.Email)
	if err != nil {
		return "", errors.New("user not found: " + err.Error())
	}

	if foundUser.Password != attempt.Password {
		return "", errors.New("wrong password")
	}

	lastLogin := time.Now()
	foundUser.LastLogin = &lastLogin
	err = userService.UserRepository.UpdateUser(foundUser)

	holderToken := fmt.Sprintf("logged in as %v, %v, %v", foundUser.Uuid, foundUser.Email, foundUser.UserType)
	return holderToken, nil
}

func (userService UserService) DeleteUser(uuid string) error {
	foundUser, err := userService.findUserByUuid(uuid)
	if err != nil {
		return err
	}
	err = userService.UserRepository.DeleteUser(foundUser)
	if err != nil {
		return err
	}
	return nil
}

func (userService UserService) findUserByEmail(email string) (*domain.User, error) {
	foundUser, err := userService.UserRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (userService UserService) findUserByUuid(uuid string) (*domain.User, error) {
	foundUser, err := userService.UserRepository.GetUserByUuid(uuid)
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
