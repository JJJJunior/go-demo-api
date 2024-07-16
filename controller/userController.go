package controller

import (
	"api/curd"
	"api/types"
	"api/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// 删除用户
func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不存在"})
		}
		if err := curd.DeleteUserByID(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "删除成功"})
	}
}

// 用户提权
func SetUserToAdmin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "id不存在"})
		}
		if err := curd.SetUserToAdmin(db, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "成功"})
	}
}

func GetUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := curd.GetUsers(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "获取成功", Data: users})
	}
}

//用户退出登录

func Logout(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := utils.GetClaimsByToken(c)
		if claims == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "token不存在"})
		}
		if err := curd.UpdateLoginStatusByID(db, claims.ID, false); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		c.Cookie(&fiber.Cookie{Name: "Authorization", Value: ""})
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "退出成功"})
	}
}

// 用户登录
func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _user types.User
		//解析json
		if err := c.BodyParser(&_user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		//校验数据格式
		validate := validator.New()
		if err := validate.Struct(_user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		//数据库逻辑
		auth, err := curd.FindAuthByUserName(db, _user.Username)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(types.Error{Error: err.Error()})
		}
		//校验密码
		result := utils.CheckPasswordHash(_user.Password, auth.HashPassword)
		if !result {
			return c.Status(fiber.StatusUnauthorized).JSON(types.Error{Error: "密码校验失败"})
		}
		//修改登录状态
		if err := curd.UpdateLoginStatusByID(db, auth.UserID, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		//查询角色
		role := curd.GetRoleByUserID(db, auth.UserID)
		if role == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "未查询到角色"})
		}
		//生成token
		token, claims := utils.GenerateJWT(&auth.UserID, &_user.Username, role)
		//生成token成功，下发token
		if token == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "生成token为空"})
		}
		curd.UpdateClaimsInAuth(db, claims)
		c.Cookie(&fiber.Cookie{Name: "Authorization", Value: fmt.Sprintf("%s +%s", "Basic ", *token)})
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "登录成功"})
	}
}

// 注册
func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var _user types.User
		//解析json
		if err := c.BodyParser(&_user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "JSON解析失败"})
		}
		//校验数据格式
		validate := validator.New()
		if err := validate.Struct(_user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}
		//操作数据库
		if curd.IsExistUserByUserName(db, _user.Username) {
			return c.Status(fiber.StatusConflict).JSON(types.Error{Error: "用户已存在"})
		}
		//生成hash密码
		hashedPassword := utils.HashPassword(_user.Password)
		if hashedPassword == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: "密码生成失败"})
		}
		err := curd.CreateOneUser(db, _user.Username, hashedPassword)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "注册成功"})
	}
}
