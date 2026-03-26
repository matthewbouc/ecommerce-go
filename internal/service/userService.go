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
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository        repository.UserRepository
	BankAccountRepository repository.BankAccountRepository
	CartRepository        repository.CartRepository
	OrderRepository       repository.OrderRepository
	Auth                  helper.Auth
	Config                config.AppConfig
	SmsClient             sms.SmsClient
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

func (userService UserService) GetProfile(userUuid uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := userService.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}
	return &dto.UserProfileResponse{
		Uuid:      user.Uuid,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		UserType:  string(user.UserType),
		Verified:  user.Verified,
	}, nil
}

func (userService UserService) UpdateProfile(userUuid uuid.UUID, input dto.UpdateProfileRequest) error {
	user, err := userService.findUserByUuid(userUuid)
	if err != nil {
		return err
	}
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Phone = input.Phone
	return userService.UserRepository.UpdateUser(user)
}

// Cart

func (userService UserService) GetCart(userUuid uuid.UUID) ([]dto.CartItemResponse, error) {
	user, err := userService.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	items, err := userService.CartRepository.GetCartByUserId(user.Id)
	if err != nil {
		return nil, err
	}

	response := make([]dto.CartItemResponse, len(items))
	for i, item := range items {
		response[i] = dto.CartItemResponse{
			Id:        item.Id,
			ProductId: item.ProductId,
			SellerId:  item.SellerId,
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			Price:     item.Price,
			Quantity:  item.Quantity,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}
	return response, nil
}

func (userService UserService) AddToCart(userUuid uuid.UUID, input dto.AddToCartRequest) error {
	user, err := userService.findUserByUuid(userUuid)
	if err != nil {
		return err
	}

	// Upsert: if this product is already in the cart, increment quantity
	existing, err := userService.CartRepository.GetCartItemByUserAndProduct(user.Id, input.ProductId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existing != nil {
		existing.Quantity += input.Quantity
		existing.Price = input.Price // refresh snapshot price
		return userService.CartRepository.UpdateCartItem(existing)
	}

	_, err = userService.CartRepository.CreateCartItem(domain.CartItem{
		UserId:    user.Id,
		ProductId: input.ProductId,
		SellerId:  input.SellerId,
		Name:      input.Name,
		ImageUrl:  input.ImageUrl,
		Price:     input.Price,
		Quantity:  input.Quantity,
	})
	return err
}

// Orders

func (userService UserService) GetOrders(userUuid uuid.UUID) ([]dto.OrderResponse, error) {
	user, err := userService.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	orders, err := userService.OrderRepository.GetOrdersByUserId(user.Id)
	if err != nil {
		return nil, err
	}

	response := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		response[i] = toOrderResponse(order)
	}
	return response, nil
}

func (userService UserService) GetOrderById(userUuid uuid.UUID, orderUuid uuid.UUID) (*dto.OrderResponse, error) {
	user, err := userService.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	order, err := userService.OrderRepository.GetOrderByUuid(user.Id, orderUuid)
	if err != nil {
		return nil, err
	}

	res := toOrderResponse(*order)
	return &res, nil
}

func (userService UserService) BecomeSeller(req dto.BecomeSellerRequest) (string, error) {

	user, err := userService.findUserByUuid(req.Uuid)
	if err != nil {
		return "", err
	}

	if user.UserType == domain.SELLER {
		return "", errors.New("user is already a seller")
	}

	user.UserType = domain.SELLER
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	// update user, create seller
	err = userService.UserRepository.UpdateUser(user)
	if err != nil {
		return "", err
	}

	_, err = userService.BankAccountRepository.CreateBankAccount(domain.BankAccount{
		UserId:            user.Id,
		BankAccountNumber: req.BankAccountNumber,
		RoutingNumber:     req.RoutingNumber,
	})

	if err != nil {
		return "", err
	}

	sellerToken, err := userService.Auth.GenerateJwt(user.Uuid, user.Email, user.UserType)
	if err != nil {
		return "", err
	}

	return sellerToken, nil
}

func toOrderResponse(order domain.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ProductId: item.ProductId,
			SellerId:  item.SellerId,
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			Price:     item.Price,
			Quantity:  item.Quantity,
		}
	}
	return dto.OrderResponse{
		Uuid:       order.Uuid,
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
		Items:      items,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}
