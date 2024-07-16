package curd

import (
	"api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 获取制定特性
func GetOnePropertyByID(db *gorm.DB, id string) *models.Property {
	var property models.Property
	if err := db.Where("id = ?", id).First(&property).Error; err != nil {
		return nil
	}
	return &property
}

func IsExistPropertyID(db *gorm.DB, id string) bool {
	var property *models.Property
	result := db.Where("id = ?", id).First(&property)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

// 获取全部产品特性
func GetAllProperties(db *gorm.DB) (*[]models.Property, error) {
	var properties []models.Property
	if err := db.Preload("Values").Order("updated_at desc").Find(&properties).Error; err != nil {
		return nil, err
	}
	return &properties, nil
}

// 添加一个特性
func CreateOneProperty(db *gorm.DB, _property *models.Property) bool {
	var property models.Property
	property.ID = uuid.NewString()
	property.Name = _property.Name
	if err := db.Create(&property).Error; err != nil {
		return false
	}
	return true
}

func SaveOneProperty(db *gorm.DB, property *models.Property) bool {
	if err := db.Save(&property).Error; err != nil {
		return false
	}
	return true
}

func DeleteOnePropertyByID(db *gorm.DB, id string) {
	var property models.Property
	if IsExistPropertyID(db, id) {
		db.Where("id = ?", id).Delete(&property)
	}
}
