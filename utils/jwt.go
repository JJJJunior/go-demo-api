package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strconv"
	"time"
)

// 定义 JWT 签名的密钥
var jwtKey = []byte(os.Getenv("JWT_KEY"))

// Claims 定义一个结构体来表示 JWT 的声明 (claims)
type Claims struct {
	UserID string
	Email  string
	Role   string
	jwt.RegisteredClaims
	UpdateTime time.Time
	Exp        time.Time
}

// GenerateJWT 根据用户名进行jwt签名
func GenerateJWT(userID *string, email *string, role *string) (*string, *Claims) {
	if userID == nil || email == nil || role == nil {
		log.Printf("GenerateJWT时发现缺少必要参数：userID和email:%s\n")
		return nil, nil
	}
	//log.Printf("调用utils/my_jwt/GenerateJWT...\n")
	// 设置 JWT的过期时间
	expirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXP_HOURS"))
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	claims := &Claims{
		UserID: *userID,
		Email:  *email,
		Role:   *role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		UpdateTime: time.Now(),
		Exp:        expirationTime,
	}
	// 创建 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// 使用密钥签名并获取完整的编码后的字符串 token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("GenerateJWT:%s\n", err)
		return nil, nil
	}
	return &tokenString, claims
}

// ValidateJWT 解析jwt
func ValidateJWT(tokenString string) *Claims {
	//log.Printf("调用utils/my_jwt/ValidateJWT...\n")
	claims := &Claims{}

	// 解析和验证 JWT
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			log.Println("invalid signature")
			return nil
		}
		log.Printf("could not parse token: %v\n", err)
		return nil
	}

	if !token.Valid {
		log.Println("ValidateJWT: invalid token")
		return nil
	}
	return claims
}
