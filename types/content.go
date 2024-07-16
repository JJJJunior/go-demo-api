package types

import "api/models"

type PaginatedSubCategories struct {
	SubCategories []models.SubCategory
	TotalPages    int
}

type PaginatedCategories struct {
	Categories []models.Category
	TotalPages int
}

type AddManyValues struct {
	NameList   []string `json:"name_list"`
	PropertyID string   `json:"property_id"`
}

type Success struct {
	Message string
	Data    interface{}
}
type Error struct {
	Error string
	Data  interface{}
}

type User struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gte=6,lte=20"`
}
