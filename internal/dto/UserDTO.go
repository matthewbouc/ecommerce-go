package dto

import "github.com/google/uuid"

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
	SwiftCode         string    `json:"swiftCode"`
	PaymentMethod     string    `json:"paymentMethod"`
}
