package initdb

import (
	"api/models"
	"api/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

func MyDB() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
		return nil
	}
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		//去除外键引用
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Error connect db")
		return nil
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Auth{},
		&models.Role{},
		&models.Product{},
		&models.Image{},
		&models.Property{},
		&models.Category{},
		&models.SubCategory{},
		&models.Value{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
		&models.Coupon{},
		&models.Customer{},
		&models.CustomerCoupon{},
		&models.ProductImage{},
		&models.ProductProperty{},
	); err != nil {
		log.Debugf("AutoMigrate: %+v\n", err)
		return nil
	}
	//初始化数据
	initSQL(db)
	return db
}

func initSQL(db *gorm.DB) {
	var category []models.Category
	result := db.Find(&category)
	if result.RowsAffected == 0 {
		db.Create(&models.Category{
			ID:   uuid.NewString(),
			Name: "default",
		})
	}
	var roles []models.Role
	if db.Find(&roles).RowsAffected == 0 {
		db.Create(&models.Role{
			Name:  "admin",
			Alias: "管理员",
		})
		db.Create(&models.Role{
			Name:  "user",
			Alias: "普通用户",
		})
	}
	var users []models.User
	var role models.Role
	db.Where("name = ?", "admin").First(&role)
	if db.Find(&users).RowsAffected == 0 {
		userID := uuid.NewString()
		user := &models.User{
			ID:       userID,
			Name:     "kelvin",
			Username: "admin",
			Email:    "admin@admin.com",
			NickName: "kelvin",
			Auth: models.Auth{
				ID:           uuid.NewString(),
				HashPassword: utils.HashPassword("sd1991308"),
				TokenExp:     time.Now(),
				UserID:       userID,
			},
			RoleID: role.ID,
		}
		db.Create(&user)
	}
}
