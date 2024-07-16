package utils

import (
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 生成密码的哈希值
func HashPassword(password string) string {
	//log.Printf("调用utils/my_password/HashPassword...\n")
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Debugf("调用utils/my_password/HashPassword:%+v\n", err)
		return ""
	}
	return string(bytes)
}

// CheckPasswordHash 验证输入的密码和哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	//log.Printf("调用utils/my_password/CheckPasswordHash...\n")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		//log.Debugf("调用utils/my_password/CheckPasswordHash:%+v\n", err)
		return false
	}
	return true
}
