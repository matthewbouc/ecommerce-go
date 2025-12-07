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

func (a *Auth) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO log the actual error
		return "", errors.New("password hashing failure")
	}
	return string(bytes), err
}

func (a *Auth) VerifyPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// TODO - Add error logging
		return errors.New("password does not match")
	}
	return nil
}

func (a *Auth) GenerateJwt(id uuid.UUID, email string, role domain.UserType) (string, error) {
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

func (a *Auth) VerifyJwt(tokenString string) (domain.User, error) {

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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.User{}, errors.New("token is expired")
		}

		strUuid := claims["userId"].(string)
		userUuid, err := uuid.Parse(strUuid)
		if err != nil {
			return domain.User{}, errors.New("invalid uuid")
		}

		userRole := domain.UserType(claims["role"].(string))
		if !userRole.IsValidUserType() {
			return domain.User{}, errors.New("invalid user type")
		}

		user := domain.User{
			Uuid:     userUuid,
			Email:    claims["email"].(string),
			UserType: userRole,
		}
		return user, nil
	}
	return domain.User{}, errors.New("token verification failed")
}

func (a *Auth) RefreshJwt(ctx fiber.Ctx) error {
	return nil
}

func (a *Auth) Authorize(ctx fiber.Ctx) error {

	authHeader := ctx.Get("Authorization", "")

	user, err := a.VerifyJwt(authHeader)

	if err == nil && user.Uuid != uuid.Nil {
		ctx.Locals("user", user)
		return ctx.Next()
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "authorization failed",
		"reason":  err,
	})
}

func (a *Auth) GetCurrentUser(ctx fiber.Ctx) domain.User {
	user := ctx.Locals("user")
	return user.(domain.User)
}

func (a *Auth) GenerateCode() (int, error) {
	return RandomNumbers(6)
}
