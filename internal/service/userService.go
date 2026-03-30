package service

import (
	"context"
	"ecommerce/config"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/pkg/notification/sms"
	catalogpb "ecommerce/proto/catalog"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Register(dto.RegisterRequest) (string, error)
	Login(dto.LoginRequest) (string, error)
	DeleteUser(uuid.UUID) error
	IsActiveUser(uuid.UUID) bool
	GetVerificationCode(domain.User) (int, error)
	VerifyCode(uuid.UUID, int) error
	GetProfile(uuid.UUID) (*dto.UserProfileResponse, error)
	UpdateProfile(uuid.UUID, dto.UpdateProfileRequest) error
	GetCart(uuid.UUID) ([]dto.CartItemResponse, error)
	AddToCart(uuid.UUID, dto.AddToCartRequest) error
	RemoveFromCart(uuid.UUID, uint) error
	GetOrders(uuid.UUID) ([]dto.OrderResponse, error)
	GetOrderById(uuid.UUID, uuid.UUID) (*dto.OrderResponse, error)
	Checkout(uuid.UUID) (*dto.OrderResponse, error)
	UpdateOrderStatus(userUuid uuid.UUID, orderUuid uuid.UUID, status string) error
	BecomeSeller(dto.BecomeSellerRequest) (string, error)
}

// userServiceImpl is the concrete implementation. Unexported intentionally —
// external packages receive a UserService interface from NewUserService.
type userServiceImpl struct {
	userRepo      repository.UserRepository
	bankRepo      repository.BankAccountRepository
	cartRepo      repository.CartRepository
	orderRepo     repository.OrderRepository
	auth          helper.Auth
	config        config.AppConfig
	smsClient     sms.SmsClient
	catalogClient catalogpb.CatalogServiceClient
}

func NewUserService(
	userRepo repository.UserRepository,
	bankRepo repository.BankAccountRepository,
	cartRepo repository.CartRepository,
	orderRepo repository.OrderRepository,
	auth helper.Auth,
	cfg config.AppConfig,
	smsClient sms.SmsClient,
	catalogClient catalogpb.CatalogServiceClient,
) UserService {
	return &userServiceImpl{
		userRepo:      userRepo,
		bankRepo:      bankRepo,
		cartRepo:      cartRepo,
		orderRepo:     orderRepo,
		auth:          auth,
		config:        cfg,
		smsClient:     smsClient,
		catalogClient: catalogClient,
	}
}

func (s *userServiceImpl) Register(userInfo dto.RegisterRequest) (string, error) {
	hashPassword, err := s.auth.HashPassword(userInfo.Password)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.CreateUser(domain.User{
		Email:     userInfo.Email,
		Password:  hashPassword,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Phone:     userInfo.Phone,
	})
	if err != nil {
		return "", err
	}

	return s.auth.GenerateJwt(user.Uuid, user.Email, user.UserType)
}

func (s *userServiceImpl) Login(attempt dto.LoginRequest) (string, error) {
	user, err := s.findUserByEmail(attempt.Email)
	if err != nil {
		return "", errors.New("user not found: " + err.Error())
	}

	err = s.auth.VerifyPassword(attempt.Password, user.Password)
	if err != nil {
		return "", err
	}

	lastLogin := time.Now()
	user.LastLogin = &lastLogin
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return s.auth.GenerateJwt(user.Uuid, user.Email, user.UserType)
}

func (s *userServiceImpl) DeleteUser(id uuid.UUID) error {
	foundUser, err := s.userRepo.GetUserByUuid(id)
	if err != nil {
		return err
	}
	return s.userRepo.DeleteUser(foundUser)
}

func (s *userServiceImpl) findUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userServiceImpl) findUserByUuid(id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetUserByUuid(id)
}

func (s *userServiceImpl) IsActiveUser(id uuid.UUID) bool {
	foundUser, err := s.userRepo.GetUserByUuid(id)
	return err == nil && !foundUser.DeletedAt.Valid
}

func (s *userServiceImpl) isVerifiedUser(id uuid.UUID) bool {
	foundUser, err := s.userRepo.GetUserByUuid(id)
	return err == nil && foundUser.Verified
}

func (s *userServiceImpl) GetVerificationCode(attempt domain.User) (int, error) {
	if s.isVerifiedUser(attempt.Uuid) {
		return 0, errors.New("user is already verified")
	}

	code, err := s.auth.GenerateCode()
	if err != nil {
		return 0, err
	}

	expiry := time.Now().Add(15 * time.Minute)
	user := domain.User{
		Uuid:             attempt.Uuid,
		Expiry:           &expiry,
		VerificationCode: code,
	}

	if err = s.userRepo.UpdateUser(&user); err != nil {
		return 0, errors.New("unable to update verification code")
	}

	msg := fmt.Sprintf("Your verification code is %v", code)
	if err = s.smsClient.SendSms(user.Phone, msg); err != nil {
		return 0, errors.New("unable to send sms message")
	}

	// TODO remove the return "code" at some point
	return code, nil
}

func (s *userServiceImpl) VerifyCode(id uuid.UUID, code int) error {
	if s.isVerifiedUser(id) {
		return errors.New("user is already verified")
	}

	user, err := s.findUserByUuid(id)
	if err != nil {
		return err
	}

	if user.VerificationCode != code {
		return errors.New("invalid verification code")
	}

	// Expiry is nil if the user never requested a code — treat as expired.
	if user.Expiry == nil || time.Now().After(*user.Expiry) {
		return errors.New("verification code is expired")
	}

	user.Verified = true
	return s.userRepo.UpdateUser(user)
}

// Profile

func (s *userServiceImpl) GetProfile(userUuid uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.findUserByUuid(userUuid)
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

func (s *userServiceImpl) UpdateProfile(userUuid uuid.UUID, input dto.UpdateProfileRequest) error {
	user, err := s.findUserByUuid(userUuid)
	if err != nil {
		return err
	}
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Phone = input.Phone
	return s.userRepo.UpdateUser(user)
}

// Cart

func (s *userServiceImpl) GetCart(userUuid uuid.UUID) ([]dto.CartItemResponse, error) {
	user, err := s.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	items, err := s.cartRepo.GetCartByUserId(user.Id)
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

func (s *userServiceImpl) RemoveFromCart(userUuid uuid.UUID, cartItemId uint) error {
	user, err := s.findUserByUuid(userUuid)
	if err != nil {
		return err
	}
	// Ownership is enforced inside the repository by scoping the delete to user.Id.
	return s.cartRepo.DeleteCartItemByIdAndUser(cartItemId, user.Id)
}

func (s *userServiceImpl) AddToCart(userUuid uuid.UUID, input dto.AddToCartRequest) error {
	user, err := s.findUserByUuid(userUuid)
	if err != nil {
		return err
	}

	// Upsert: if this product is already in the cart, increment quantity.
	existing, err := s.cartRepo.GetCartItemByUserAndProduct(user.Id, input.ProductId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existing != nil {
		existing.Quantity += input.Quantity
		existing.Price = input.Price // refresh snapshot price
		return s.cartRepo.UpdateCartItem(existing)
	}

	_, err = s.cartRepo.CreateCartItem(domain.CartItem{
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

func (s *userServiceImpl) GetOrders(userUuid uuid.UUID) ([]dto.OrderResponse, error) {
	user, err := s.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.GetOrdersByUserId(user.Id)
	if err != nil {
		return nil, err
	}

	response := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		response[i] = toOrderResponse(order)
	}
	return response, nil
}

func (s *userServiceImpl) GetOrderById(userUuid uuid.UUID, orderUuid uuid.UUID) (*dto.OrderResponse, error) {
	user, err := s.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.GetOrderByUuid(user.Id, orderUuid)
	if err != nil {
		return nil, err
	}

	res := toOrderResponse(*order)
	return &res, nil
}

func (s *userServiceImpl) Checkout(userUuid uuid.UUID) (*dto.OrderResponse, error) {
	user, err := s.findUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	cartItems, err := s.cartRepo.GetCartByUserId(user.Id)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Give the gRPC calls a 10s budget collectively.
	// TODO: propagate the HTTP request context through service methods instead of creating a fresh background context here.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var orderItems []domain.OrderItem
	var totalPrice float64

	for _, cartItem := range cartItems {
		// Call the Catalog gRPC server — even though it runs in the same process,
		// this goes through the full gRPC stack (serialize → TCP loopback → deserialize).
		// When Order Service splits into its own binary, only the address will need to change.
		resp, err := s.catalogClient.GetProduct(ctx, &catalogpb.GetProductRequest{
			Uuid: cartItem.ProductId.String(),
		})
		if err != nil {
			return nil, fmt.Errorf("product %s is no longer available", cartItem.Name)
		}

		product := resp.Product
		if int(product.Stock) < int(cartItem.Quantity) {
			return nil, fmt.Errorf("insufficient stock for '%s': requested %d, available %d",
				product.Name, cartItem.Quantity, product.Stock)
		}

		// Use the authoritative price from Catalog, not the snapshot stored in
		// the cart. This catches any price changes since the item was added.
		totalPrice += product.Price * float64(cartItem.Quantity)

		orderItems = append(orderItems, domain.OrderItem{
			ProductId: cartItem.ProductId,
			SellerId:  cartItem.SellerId,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Price:     product.Price,
			Quantity:  cartItem.Quantity,
		})
	}

	order, err := s.orderRepo.CreateOrder(domain.Order{
		UserId:     user.Id,
		Status:     domain.OrderPending,
		TotalPrice: totalPrice,
		Items:      orderItems,
	})
	if err != nil {
		return nil, err
	}

	// TODO: wrap CreateOrder + ClearCartByUserId in a single DB transaction
	// so a cart-clear failure doesn't leave the user with a stale cart.
	if err := s.cartRepo.ClearCartByUserId(user.Id); err != nil {
		fmt.Printf("[checkout] warning: order %s created but cart not cleared: %v\n", order.Uuid, err)
	}

	res := toOrderResponse(order)
	return &res, nil
}

func (s *userServiceImpl) UpdateOrderStatus(userUuid uuid.UUID, orderUuid uuid.UUID, status string) error {
	orderStatus := domain.OrderStatus(status)
	if !orderStatus.IsValidStatus() {
		return fmt.Errorf("invalid status '%s'", status)
	}
	// TODO: verify userUuid has permission to update this order
	// (buyer can cancel pending orders; seller can advance to shipped/delivered)
	return s.orderRepo.UpdateOrderStatus(orderUuid, orderStatus)
}

func (s *userServiceImpl) BecomeSeller(req dto.BecomeSellerRequest) (string, error) {
	user, err := s.findUserByUuid(req.Uuid)
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

	if err = s.userRepo.UpdateUser(user); err != nil {
		return "", err
	}

	_, err = s.bankRepo.CreateBankAccount(domain.BankAccount{
		UserId:            user.Id,
		BankAccountNumber: req.BankAccountNumber,
		RoutingNumber:     req.RoutingNumber,
	})
	if err != nil {
		return "", err
	}

	return s.auth.GenerateJwt(user.Uuid, user.Email, user.UserType)
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
