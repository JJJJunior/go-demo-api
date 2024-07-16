package curd

import (
	"api/models"
	"api/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
)

func GetOneSubCategoryByID(db *gorm.DB, id *string) *models.SubCategory {
	var subCategory models.SubCategory
	if err := db.Where("id = ?", id).First(&subCategory).Error; err != nil {
		return nil
	}
	return &subCategory
}

// 更新一个二级栏目
func UpdateOneSubCategory(db *gorm.DB, subCategory models.SubCategory) error {
	if err := db.Save(&subCategory).Error; err != nil {
		return err
	}
	return nil
}

// 新增栏目
func CreateSubCategory(db *gorm.DB, subCategory models.SubCategory) error {
	subCategory.ID = uuid.NewString()
	if err := db.Create(&subCategory).Error; err != nil {
		return err
	}
	return nil
}

// 二级栏目的分页查询
func GetSubCategoriesByPaginate(db *gorm.DB, page int, pageSize int) (*types.PaginatedSubCategories, error) {
	var subCategories []models.SubCategory
	var totalRecords int64

	// 获取满足条件的总记录数
	if err := db.Model(models.SubCategory{}).Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	offset := (page - 1) * pageSize
	if err := db.Preload("Properties.Values").Order("updated_at desc").Limit(pageSize).Offset(offset).Find(&subCategories).Error; err != nil {
		return nil, err
	}
	return &types.PaginatedSubCategories{
		SubCategories: subCategories,
		TotalPages:    totalPages,
	}, nil
}

// 根据ID删除栏目
func DeleteOneSubCategoryByID(db *gorm.DB, subCategoryID string) error {
	var subCategory models.SubCategory
	if err := db.Where("id = ?", subCategoryID).Delete(&subCategory).Error; err != nil {
		return err
	}
	return nil
}
