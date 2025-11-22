package service

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepository repository.UserRepository
	Auth           helper.Auth
}

func (userService UserService) Register(userInfo dto.RegisterDTO) (string, error) {

	hashPassword, err := userService.Auth.HashPassword(userInfo.Password)
	if err != nil {
		return "", err
	}

	user, err := userService.UserRepository.CreateUser(domain.User{
		Email:     userInfo.Email,
		Password:  hashPassword,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Phone:     userInfo.Phone,
	})
	if err != nil {
		return "", err
	}

	return userService.Auth.GenerateJwt(user.Uuid, user.Email, user.UserType)
}

func (userService UserService) Login(attempt dto.LoginDTO) (string, error) {
	user, err := userService.findUserByEmail(attempt.Email)
	if err != nil {
		return "", errors.New("user not found: " + err.Error())
	}

	err = userService.Auth.VerifyPassword(attempt.Password, user.Password)
	if err != nil {
		return "", err
	}

	lastLogin := time.Now()
	user.LastLogin = &lastLogin
	err = userService.UserRepository.UpdateUser(user)

	return userService.Auth.GenerateJwt(user.Uuid, user.Email, user.UserType)
}

func (userService UserService) DeleteUser(uuid uuid.UUID) error {
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

func (userService UserService) findUserByUuid(uuid uuid.UUID) (*domain.User, error) {
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
