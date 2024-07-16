package curd

import (
	"api/models"
	"api/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
)

func GetAllFromCategoryToValue(db *gorm.DB) (*[]models.Category, error) {
	var categories *[]models.Category
	if err := db.Preload("SubCategories.Properties.Values").Order("updated_at desc").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// 删除一个一级栏目
func DeleteOneCategory(db *gorm.DB, id string) bool {
	var category models.Category
	if err := db.Where("id = ?", id).Delete(&category).Error; err != nil {
		return false
	}
	return true
}

// 创建一个一级栏目
func CreateOneCategory(db *gorm.DB, _category models.Category) error {
	if err := db.Create(&models.Category{
		ID:   uuid.NewString(),
		Name: _category.Name,
	}).Error; err != nil {
		return err
	}
	return nil
}

// 一级栏目的分页查询
func GetCategoriesByPaginate(db *gorm.DB, page int, pageSize int) (*types.PaginatedCategories, error) {
	var categories []models.Category
	var totalRecords int64
	// 获取满足条件的总记录数
	if err := db.Model(models.Category{}).Count(&totalRecords).Error; err != nil {
		return nil, err
	}
	// 计算总页数
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	offset := (page - 1) * pageSize
	if err := db.Preload("SubCategories.Properties.Values").Order("updated_at desc").Limit(pageSize).Offset(offset).Find(&categories).Error; err != nil {
		return nil, err
	}
	return &types.PaginatedCategories{
		Categories: categories,
		TotalPages: totalPages,
	}, nil
}

// 获取默认parent的id
func GetDefaultCategoryID(db *gorm.DB) (*models.Category, error) {
	var category models.Category
	if err := db.Where("name = ?", "default").First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// 根据一级栏目ID查找一级栏目
func FindCategoryByCategoryID(db *gorm.DB, categoryID string) (*models.Category, *error) {
	var category models.Category
	if err := db.Where("id = ?", categoryID).First(&category).Error; err != nil {
		return nil, &err
	}
	return &category, nil
}
