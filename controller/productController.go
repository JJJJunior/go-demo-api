package controller

import (
	"api/curd"
	"api/models"
	"api/types"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var product models.Product
		if err := c.BodyParser(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if err := curd.CreateProduct(db, &product); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "创建成功"})
	}
}

func GetProducts(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var products []models.Product
		if err := curd.GetProducts(db, &products); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: products})
	}
}

func GetProduct(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var product models.Product
		if err := curd.GetProduct(db, id, &product); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(types.Error{Error: "Product not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: product})
	}
}

func UpdateProduct(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var product models.Product
		if err := c.BodyParser(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if err := curd.UpdateProduct(db, &product); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "修改成功"})
	}
}

func DeleteProduct(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := curd.DeleteProduct(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "删除失败"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "删除成功"})
	}
}
