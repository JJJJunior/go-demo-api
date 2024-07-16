package controller

import (
	"api/curd"
	"api/models"
	"api/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOneValueByID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不存在"})
		}
		value, err := curd.GetOneValueByID(db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: value})
	}
}

func CreateManyValues(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _addManyValues types.AddManyValues
		//解析json
		if err := c.BodyParser(&_addManyValues); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		if !curd.CreateManyValues(db, _addManyValues) {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "添加失败"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "创建成功",
			Data:    nil,
		})
	}
}

func CreateOneValue(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _propertyValue models.Value
		//解析json
		if err := c.BodyParser(&_propertyValue); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		//校验数据格式
		validate := validator.New()
		if err := validate.Struct(_propertyValue); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if err := curd.SaveOneValue(db, &_propertyValue); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "创建成功",
			Data:    nil,
		})
	}
}

func UpdateOneValue(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var value models.Value
		//解析json
		if err := c.BodyParser(&value); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		//校验数据格式
		validate := validator.New()
		if err := validate.Struct(value); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if err := curd.SaveOneValue(db, &value); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "创建成功",
			Data:    nil,
		})
	}
}

func DeleteOneValue(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不能为空"})
		}
		if err := curd.DeleteOneValueByID(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "删除成功",
			Data:    nil,
		})
	}
}

func GetAllValues(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		values, err := curd.GetAllValues(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "创建成功",
			Data:    values,
		})
	}
}
