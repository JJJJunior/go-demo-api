package middleware

import (
	"api/types"
	"api/utils"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//log.Printf("AuthMiddleware...\n")
		cookie := c.Cookies("Authorization")
		if cookie != "" {
			//spew.Dump(cookie)
			token := cookie[strings.Index(cookie, "eyJ"):]
			//spew.Dump(token)
			claims := utils.ValidateJWT(token)
			//spew.Dump(claims)
			if claims != nil {
				if claims.Role == "管理员" {
					err := c.Status(fiber.StatusOK).JSON(types.Success{Message: "鉴权成功"})
					if err != nil {
						return err
					}
					return c.Next()
				} else {
					return c.Status(fiber.StatusUnauthorized).JSON(types.Error{Error: "不是管理员"})
				}
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(types.Error{Error: "无权限"})
			}
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(types.Error{Error: "无权限"})
		}
	}
}
