package service

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"log"
)

type UserService struct {
}

// func (serv UserService) SignUp(input any) {} // any and interface don't enforce types, beneficial if it's a changing input
// func (serv UserService) SignUp(input interface) {}
func (serv UserService) Register(userInfo dto.RegisterDTO) (string, error) {
	log.Println(userInfo)

	// TODO: setup db to create user here
	return "this-is-the-token", nil
}

func (serv UserService) Login(attempt dto.LoginDTO) (string, error) {
	return "", nil
}

func (serv UserService) findUserByEmail(email string) (*domain.User, error) {
	return nil, nil // using nil, nil because *domain is using a pointer to deal directly with the User object
}

func (serv UserService) GetVerificationCode(attempt domain.User) (int, error) {
	return 0, nil
}

func (serv UserService) VerifyCode(id uint, code int) error {
	return nil
}

func (serv UserService) CreateProfile(id uint, input any) error {
	return nil
}

func (serv UserService) GetProfile(id uint) (*domain.User, error) {
	return nil, nil
}

func (serv UserService) UpdateProfile(id uint, input any) error {
	return nil
}

func (serv UserService) BecomeSeller(id uint, input any) error {
	return nil
}

func (serv UserService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}

func (serv UserService) CreateCart(id uint, input any) error {
	return nil
}

func (serv UserService) CreateOrder(id uint, input any) error {
	return nil
}

func (serv UserService) GetOrders(id uint, input any) error {
	return nil
}

func (serv UserService) GetOrderById(id uint, input any) error {
	return nil
}
