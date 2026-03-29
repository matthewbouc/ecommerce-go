package helper

import (
	"ecommerce/internal/domain"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetupAuth(secret string) Auth {
	return Auth{
		Secret: secret,
	}
}

func (a Auth) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO log the actual error
		return "", errors.New("password hashing failure")
	}
	return string(bytes), err
}

func (a Auth) VerifyPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// TODO - Add error logging
		return errors.New("password does not match")
	}
	return nil
}

func (a Auth) GenerateJwt(id uuid.UUID, email string, role domain.UserType) (string, error) {
	if id == uuid.Nil || email == "" || role == "" {
		return "", errors.New("required inputs are missing to generate tokens")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"email":  email,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		// TODO - log error
		return "", errors.New("error generating jwt")
	}

	return tokenStr, nil
}

func (a Auth) VerifyJwt(tokenString string) (domain.User, error) {

	tokenArray := strings.Split(tokenString, " ")
	if len(tokenArray) != 2 || tokenArray[0] != "Bearer" {
		return domain.User{}, errors.New("invalid token")
	}

	token, err := jwt.Parse(tokenArray[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		return domain.User{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain.User{}, errors.New("token verification failed")
	}

	// Use comma-ok on every claim — a missing or wrong-typed claim must
	// return an error, not panic.
	strUuid, ok := claims["userId"].(string)
	if !ok || strUuid == "" {
		return domain.User{}, errors.New("invalid token claims")
	}
	userUuid, err := uuid.Parse(strUuid)
	if err != nil {
		return domain.User{}, errors.New("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return domain.User{}, errors.New("invalid token claims")
	}

	roleStr, ok := claims["role"].(string)
	if !ok || roleStr == "" {
		return domain.User{}, errors.New("invalid token claims")
	}
	userRole := domain.UserType(roleStr)
	if !userRole.IsValidUserType() {
		return domain.User{}, errors.New("invalid user role")
	}

	return domain.User{
		Uuid:     userUuid,
		Email:    email,
		UserType: userRole,
	}, nil
}

func (a Auth) RefreshJwt(ctx fiber.Ctx) error {
	return nil
}

func (a Auth) RequireRole(role domain.UserType) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		user := a.GetCurrentUser(ctx)
		if user.UserType != role {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "forbidden: insufficient role",
			})
		}
		return ctx.Next()
	}
}

func (a Auth) Authorize(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization", "")

	user, err := a.VerifyJwt(authHeader)
	if err != nil || user.Uuid == uuid.Nil {
		// TODO: Log internally (replace with structured logger when available).
		fmt.Printf("[auth] authorization failed: %v\n", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	ctx.Locals("user", user)
	return ctx.Next()
}

func (a Auth) GetCurrentUser(ctx fiber.Ctx) domain.User {
	user, ok := ctx.Locals("user").(domain.User)
	if !ok {
		panic("GetCurrentUser: user not found in context; ensure Authorize middleware is applied to this route")
	}
	return user
}

func (a Auth) GenerateCode() (int, error) {
	return RandomNumbers(6)
}
