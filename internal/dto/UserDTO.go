package dto

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email     string `json:"email"`    // required
	Password  string `json:"password"` //required
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerificationCodeInput struct {
	Code int `json:"code"`
}

type BecomeSellerRequest struct {
	Uuid              uuid.UUID `json:"uuid,omitempty"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Phone             string    `json:"phone"`
	BankAccountNumber uint      `json:"bankAccountNumber"`
	RoutingNumber     uint      `json:"routingNumber"`
	SwiftCode         uint      `json:"swiftCode"`
	PaymentMethod     string    `json:"paymentMethod"`
}

func (r *BecomeSellerRequest) ValidateBecomeSellerRequest() error {
	hasRoutingNumber := r.RoutingNumber != 0
	hasSwiftCode := r.SwiftCode != 0

	if hasRoutingNumber && hasSwiftCode {
		return errors.New("cannot provide both routing number and swift code")
	}
	if !hasRoutingNumber && !hasSwiftCode {
		return errors.New("must provide either routing number or swift code")
	}
	return nil
}

// Profile

type UserProfileResponse struct {
	Uuid      uuid.UUID `json:"uuid"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	UserType  string    `json:"userType"`
	Verified  bool      `json:"verified"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

// Cart

type AddToCartRequest struct {
	ProductId uuid.UUID `json:"productId"`
	SellerId  uuid.UUID `json:"sellerId"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"imageUrl"`
	Price     float64   `json:"price"`
	Quantity  uint      `json:"quantity"`
}

type CartItemResponse struct {
	Id        uint      `json:"id"`
	ProductId uuid.UUID `json:"productId"`
	SellerId  uuid.UUID `json:"sellerId"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"imageUrl"`
	Price     float64   `json:"price"`
	Quantity  uint      `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Orders

type OrderItemResponse struct {
	ProductId uuid.UUID `json:"productId"`
	SellerId  uuid.UUID `json:"sellerId"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"imageUrl"`
	Price     float64   `json:"price"`
	Quantity  uint      `json:"quantity"`
}

type OrderResponse struct {
	Uuid       uuid.UUID           `json:"uuid"`
	Status     string              `json:"status"`
	TotalPrice float64             `json:"totalPrice"`
	Items      []OrderItemResponse `json:"items"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
}
