package controller

import (
	"api/curd"
	"api/models"
	"api/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func GetAllFromCategoryToValue(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		categories, err := curd.GetAllFromCategoryToValue(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: categories})
	}
}

// 删除一个一级栏目
func DeleteOneCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不存在"})
		}
		if !curd.DeleteOneCategory(db, id) {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "删除失败"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "删除成功"})
	}
}

// 创建一个一级栏目
func CreateOneCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _category models.Category
		//解析json
		if err := c.BodyParser(&_category); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		//校验数据格式
		validate := validator.New()
		if err := validate.Struct(_category); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if err := curd.CreateOneCategory(db, _category); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "创建成功",
			Data:    nil,
		})
	}
}

// 一级栏目分页查询
func GetCategories(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var page int
		var pageSize int
		pageStr := c.Query("page")
		pageSizeStr := c.Query("pageSize")

		if pageStr == "" {
			pageStr = "1"
		}
		if pageSizeStr == "" {
			pageSizeStr = "9999"
		}
		page, _ = strconv.Atoi(pageStr)
		pageSize, _ = strconv.Atoi(pageSizeStr)
		categories, err := curd.GetCategoriesByPaginate(db, page, pageSize)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: categories})
	}
}
