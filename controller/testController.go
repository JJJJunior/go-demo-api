package controller

import (
	"api/types"
	"github.com/gofiber/fiber/v2"
)

func Test() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("Authorization")
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "访问成功", Data: cookie})
	}
}
