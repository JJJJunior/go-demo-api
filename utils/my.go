package utils

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetClaimsByToken(c *fiber.Ctx) *Claims {
	cookie := c.Cookies("Authorization")
	if cookie == "" {
		return nil
	}
	token := cookie[strings.Index(cookie, "eyJ"):]
	claims := ValidateJWT(token)
	return claims
}
