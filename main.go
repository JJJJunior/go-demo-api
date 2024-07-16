package main

import (
	"api/controller"
	"api/initdb"
	"api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db := initdb.MyDB()
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:       "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
	}))
	// 配置 CORS 中间件
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://175.178.190.62:3000",             // 允许的本地主机
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",                                  // 允许的方法
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization,x-requested-with", // 允许的请求头
		AllowCredentials: true,                                                           // 允许凭据
	}))

	api := app.Group("/api")
	api.Post("/login", controller.Login(db))
	api.Get("/logout", controller.Logout(db))
	api.Post("/register", middleware.AuthMiddleware(), controller.Register(db))
	//获取值
	api.Get("/category/values", middleware.AuthMiddleware(), controller.GetAllValues(db))
	api.Get("/category/values/:id", middleware.AuthMiddleware(), controller.GetOneValueByID(db))
	api.Post("/category/values", middleware.AuthMiddleware(), controller.CreateOneValue(db))
	api.Post("/category/manyValues", middleware.AuthMiddleware(), controller.CreateManyValues(db))
	// /category/values?id=?
	api.Put("/category/values", middleware.AuthMiddleware(), controller.UpdateOneValue(db))
	// /category/values?id=?
	api.Delete("/category/values/:id", middleware.AuthMiddleware(), controller.DeleteOneValue(db))

	//获取特性
	api.Get("/category/properties", middleware.AuthMiddleware(), controller.GetAllProperties(db))
	api.Get("/category/properties/:id", middleware.AuthMiddleware(), controller.GetOnePropertyByID(db))
	api.Post("/category/properties", middleware.AuthMiddleware(), controller.CreateOneProperty(db))
	api.Put("/category/properties", middleware.AuthMiddleware(), controller.UpdateOneProperty(db))
	api.Delete("/category/properties/:id", middleware.AuthMiddleware(), controller.DeleteOneProperty(db))

	//获取子栏目 subCategories?page=?&pageSize=?
	api.Get("/category/subCategories", middleware.AuthMiddleware(), controller.GetSubCategories(db))
	api.Get("/category/subCategories/:id", middleware.AuthMiddleware(), controller.GetOneSubCategoryByID(db))
	api.Post("/category/subCategories", middleware.AuthMiddleware(), controller.CreateSubCategoryAndAll(db))
	api.Put("/category/subCategories", middleware.AuthMiddleware(), controller.UpdateOneSubCategory(db))
	api.Delete("/category/subCategories/:id", middleware.AuthMiddleware(), controller.DeleteOneSubCategory(db))

	//获取父栏目
	// categories?id=?
	api.Delete("/category/categories", middleware.AuthMiddleware(), controller.DeleteOneCategory(db))
	// categories?page=?&pageSize=?
	api.Get("/category/categories", middleware.AuthMiddleware(), controller.GetCategories(db))
	api.Post("/category/categories", middleware.AuthMiddleware(), controller.CreateOneCategory(db))
	api.Get("/category/categories/all", middleware.AuthMiddleware(), controller.GetAllFromCategoryToValue(db))

	//图片上传接口
	api.Post("/upload", middleware.AuthMiddleware(), controller.UploadPic(db))

	//产品增删改查
	api.Post("/products", middleware.AuthMiddleware(), controller.CreateProduct(db))
	api.Get("/products", middleware.AuthMiddleware(), controller.GetProducts(db))
	api.Get("/products/:id", middleware.AuthMiddleware(), controller.GetProduct(db))
	api.Put("/products/", middleware.AuthMiddleware(), controller.UpdateProduct(db))
	api.Delete("/products/:id", middleware.AuthMiddleware(), controller.DeleteProduct(db))

	api.Get("/test", middleware.AuthMiddleware(), controller.Test())
	api.Get("/users", middleware.AuthMiddleware(), controller.GetUsers(db))
	api.Delete("/users/:id", middleware.AuthMiddleware(), controller.DeleteUser(db))
	api.Get("/users/set_admin/:id", middleware.AuthMiddleware(), controller.SetUserToAdmin(db))
	app.Listen(":8000")
}
