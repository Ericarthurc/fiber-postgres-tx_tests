package middleware

import (
	"eric/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthenticationMiddleware(c *fiber.Ctx) error {
	cookieToken := c.Cookies("authToken")
	if cookieToken == "" {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": "authentication failed, no token"})
	}

	err := utils.VerifyJWTToken(c, cookieToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	return c.Next()
}
