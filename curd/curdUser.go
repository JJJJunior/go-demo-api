package curd

import (
	"api/models"
	"api/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func DeleteUserByID(db *gorm.DB, id string) error {
	var user models.User
	if err := db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func SetUserToAdmin(db *gorm.DB, id string) error {
	var role models.Role
	if err := db.Where("name = ?", "admin").First(&role).Error; err != nil {
		return err
	}
	var user models.User
	if err := db.Model(&user).Where("id = ?", id).Update("role_id", role.ID).Error; err != nil {
		return err
	}
	return nil
}

// 根据用户ID查角色
func GetRoleByUserID(db *gorm.DB, id string) *string {
	var user models.User
	if err := db.Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return nil
	}
	return &user.Role.Alias
}

func GetUsers(db *gorm.DB) (*[]models.User, error) {
	var users []models.User
	if err := db.Preload("Role").Preload("Auth").Order("updated_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

// 在数据库更新一个token失效时间
func UpdateClaimsInAuth(db *gorm.DB, claims *utils.Claims) {
	var auth models.Auth
	result := db.Model(&auth).Where("user_id = ?", claims.UserID).Update("token_exp", claims.Exp)
	if result.Error != nil {
		log.Debugf("UpdateClaimsInAuth: %+v\n", result.Error)
	}
}

// 修改用户登录状态
func UpdateLoginStatusByID(db *gorm.DB, userID string, value bool) error {
	var auth models.Auth
	if err := db.Model(&auth).Where("user_id = ?", userID).Update("is_login", value).Error; err != nil {
		return err
	}
	return nil
}

// 通过ID查权限
func FindRoleByID(db *gorm.DB, userID string) (*models.Role, error) {
	var role models.Role
	if err := db.Preload("Users").First(&role, userID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// 通过username查询密码
func FindAuthByUserName(db *gorm.DB, username string) (*models.Auth, error) {
	var user models.User
	if err := db.Where("username = ?", username).Preload("Auth").First(&user).Error; err != nil {
		return nil, err
	}
	return &user.Auth, nil
}

// 通过username查询用户是否存在
func IsExistUserByUserName(db *gorm.DB, username string) bool {
	var user models.User
	count := db.Where("username = ?", username).First(&user).RowsAffected
	if count > 0 {
		return true
	}
	return false
}

// 创建一个用户
func CreateOneUser(db *gorm.DB, username string, password string) error {
	userID := uuid.New().String()
	if err := db.Create(&models.User{
		ID:       userID,
		Email:    "",
		Username: username,
		Picture:  "",
		Name:     "",
		NickName: "",
		Auth: models.Auth{
			HashPassword: password,
			TokenExp:     time.Now(),
			UserID:       userID,
		},
	}).Error; err != nil {
		return err
	}
	return nil
}
