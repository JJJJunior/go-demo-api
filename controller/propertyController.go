package controller

import (
	"api/curd"
	"api/models"
	"api/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOnePropertyByID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不存在"})
		}
		property := curd.GetOnePropertyByID(db, id)
		if property == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "数据查询失败"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: property})
	}
}

func GetAllProperties(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		values, err := curd.GetAllProperties(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "获取成功",
			Data:    values,
		})
	}
}

func CreateOneProperty(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _property models.Property
		//解析json
		if err := c.BodyParser(&_property); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		//校验数据
		validate := validator.New()
		if err := validate.Struct(&_property); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if !curd.CreateOneProperty(db, &_property) {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "创建失败"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "创建成功",
			Data:    "",
		})
	}
}

func UpdateOneProperty(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _property models.Property
		//解析json
		if err := c.BodyParser(&_property); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		//校验数据
		validate := validator.New()
		if err := validate.Struct(&_property); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if !curd.SaveOneProperty(db, &_property) {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "更新失败"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "更新成功",
		})
	}
}

func DeleteOneProperty(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不能为空"})
		}
		curd.DeleteOnePropertyByID(db, id)
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "删除成功",
		})
	}
}
