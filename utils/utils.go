package utils

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func SendAuthCookie(c *fiber.Ctx, token string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "authToken"
	cookie.Value = token
	cookie.Expires = time.Now().Add(2 * time.Hour)
	cookie.HTTPOnly = true
	if os.Getenv("ENVIROMENT") == "production" {
		cookie.Secure = true
	}

	c.Cookie(cookie)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func CreateJWTToken(id string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyJWTToken(c *fiber.Ctx, cookieToken string) error {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookieToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return errors.New("authentication failed")
	}

	c.Locals("userID", claims["id"])
	c.Locals("authToken", cookieToken)

	return nil
}
