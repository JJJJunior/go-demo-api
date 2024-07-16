package controller

import (
	"api/curd"
	"api/models"
	"api/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// 获取一个二级栏目
func GetOneSubCategoryByID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不存在"})
		}
		subCategory := curd.GetOneSubCategoryByID(db, &id)
		if subCategory == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "子栏目不存在"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: subCategory})
	}
}

// 更新一个二级栏目
func UpdateOneSubCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _subCategory models.SubCategory
		//解析json
		if err := c.BodyParser(&_subCategory); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		//校验数据格式
		validate := validator.New()
		if err := validate.Struct(_subCategory); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if err := curd.UpdateOneSubCategory(db, _subCategory); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "更新成功"})
	}
}

// 删除一个二级栏目
func DeleteOneSubCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id参数不存在"})
		}
		if err := curd.DeleteOneSubCategoryByID(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "删除数据失败"})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "删除成功"})
	}
}

// 二级栏目分页查询
func GetSubCategories(db *gorm.DB) fiber.Handler {
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
		subCategories, err := curd.GetSubCategoriesByPaginate(db, page, pageSize)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: subCategories})
	}
}

// 新增栏目
func CreateSubCategoryAndAll(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _subCategory models.SubCategory
		//解析json
		if err := c.BodyParser(&_subCategory); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		//spew.Dump(_category)
		//校验数据格式
		validate := validator.New()
		if err := validate.Struct(_subCategory); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		if err := curd.CreateSubCategory(db, _subCategory); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "新增栏目成功"})
	}
}

// 删除栏目
func DeleteSubCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := strings.Trim(" ", c.Query("id"))
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "query参数不正确"})
		}
		err := curd.DeleteOneSubCategoryByID(db, id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{
			Message: "删除栏目成功",
			Data:    nil,
		})
	}
}
