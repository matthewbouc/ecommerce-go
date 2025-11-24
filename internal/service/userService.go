package service

import (
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/pkg/notification/sms"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepository repository.UserRepository
	Auth           helper.Auth
	Config         config.AppConfig
	SmsClient      sms.SmsClient
}

func (userService UserService) Register(userInfo dto.RegisterRequest) (string, error) {

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

func (userService UserService) Login(attempt dto.LoginRequest) (string, error) {
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
	foundUser, err := userService.UserRepository.GetUserByUuid(uuid)
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

func (userService UserService) isActiveUser(uuid uuid.UUID) bool {
	foundUser, err := userService.UserRepository.GetUserByUuid(uuid)
	return err == nil && !foundUser.DeletedAt.Valid
}

func (userService UserService) isVerifiedUser(id uuid.UUID) bool {
	foundUser, err := userService.UserRepository.GetUserByUuid(id)
	return err == nil && foundUser.Verified

}

func (userService UserService) GetVerificationCode(attempt domain.User) (int, error) {
	if userService.isVerifiedUser(attempt.Uuid) {
		return 0, errors.New("user is already verified")
	}

	code, err := userService.Auth.GenerateCode()
	if err != nil {
		return 0, err
	}

	expiry := time.Now().Add(15 * time.Minute)
	user := domain.User{
		Uuid:             attempt.Uuid,
		Expiry:           &expiry,
		VerificationCode: code,
	}

	err = userService.UserRepository.UpdateUser(&user)

	if err != nil {
		return 0, errors.New("unable to updated verification code")
	}

	// Send SMS Notification
	msg := fmt.Sprintf("Your verification code is %v", code)

	err = userService.SmsClient.SendSms(user.Phone, msg)
	if err != nil {
		return 0, errors.New("unable to send sms message")
	}

	// TODO remove the return "code" at some point
	return code, nil
}

func (userService UserService) VerifyCode(Uuid uuid.UUID, code int) error {

	if userService.isVerifiedUser(Uuid) {
		return errors.New("user is already verified")
	}

	user, err := userService.findUserByUuid(Uuid)
	if err != nil {
		return err
	}

	if user.VerificationCode != code {
		return errors.New("invalid verification code")
	}

	if time.Now().After(*user.Expiry) {
		return errors.New("verification code is expired")
	}

	user.Verified = true

	err = userService.UserRepository.UpdateUser(user)
	if err != nil {
		return err
	}

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
