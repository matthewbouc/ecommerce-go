package dto

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
