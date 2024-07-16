package curd

import (
	"api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func UpdateProduct(db *gorm.DB, product *models.Product) error {
	product.UpdatedAt = time.Now()
	product.CreatedAt = time.Now()
	var existingProduct models.Product
	tx := db.Begin()

	// 查找产品
	if err := tx.Preload("SubCategories").Preload("Properties").Preload("Images").Where("id = ?", product.ID).First(&existingProduct).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 清理旧的关联关系
	if err := tx.Model(&existingProduct).Association("SubCategories").Clear(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&existingProduct).Association("Properties").Clear(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&existingProduct).Association("Images").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 更新产品信息
	if err := tx.Model(&existingProduct).Save(product).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 重新创建新的关联关系
	if err := tx.Model(&existingProduct).Association("SubCategories").Replace(product.SubCategories); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&existingProduct).Association("Properties").Replace(product.Properties); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&existingProduct).Association("Images").Replace(product.Images); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func CreateProduct(db *gorm.DB, product *models.Product) error {
	product.ID = uuid.NewString()
	return db.Create(&product).Error
}

func GetProducts(db *gorm.DB, products *[]models.Product) error {
	return db.Preload("SubCategories").Preload("Properties").Preload("Images").Find(&products).Error
}

func GetProduct(db *gorm.DB, id string, product *models.Product) error {
	return db.Preload("SubCategories").Preload("Properties").Preload("Images").Where("id = ?", id).First(&product).Error
}

func DeleteProduct(db *gorm.DB, id string) error {
	var product models.Product
	tx := db.Begin()

	// 使用事务对象 tx 进行查询和操作
	if err := tx.Preload("SubCategories").Preload("Properties").Preload("Images").Where("id = ?", id).First(&product).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 清理关联关系
	if err := tx.Model(&product).Association("SubCategories").Clear(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&product).Association("Properties").Clear(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&product).Association("Images").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 删除产品
	if err := tx.Delete(&product).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
