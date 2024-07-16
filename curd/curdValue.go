package curd

import (
	"api/models"
	"api/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetOneValueByID(db *gorm.DB, id string) (*models.Value, error) {
	var value models.Value
	if err := db.Where("id = ?", id).First(&value).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

// 添加多个value
func CreateManyValues(db *gorm.DB, _addManyValues types.AddManyValues) bool {
	if _addManyValues.PropertyID != "" && !IsExistPropertyID(db, _addManyValues.PropertyID) {
		return false
	}
	for _, name := range _addManyValues.NameList {
		db.Create(&models.Value{
			ID:         uuid.NewString(),
			Name:       name,
			PropertyID: &_addManyValues.PropertyID,
		})
	}
	return true
}

func IsExistValueID(db *gorm.DB, id string) bool {
	var value *models.Value
	result := db.Where("id = ?", id).First(&value)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

// 新增或者更新操作，传id就更新，不传就新建
func SaveOneValue(db *gorm.DB, value *models.Value) error {
	spew.Dump(value)
	if err := db.Save(&value).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOneValueByID(db *gorm.DB, id string) error {
	var value models.Value
	if err := db.Where("id = ?", id).Delete(&value).Error; err != nil {
		return err
	}
	return nil
}

func GetAllValues(db *gorm.DB) (*[]models.Value, error) {
	var values []models.Value
	if err := db.Order("updated_at desc").Find(&values).Error; err != nil {
		return nil, err
	}
	return &values, nil
}
