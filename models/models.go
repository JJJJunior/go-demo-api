package models

import (
	"gorm.io/gorm"
	"time"
)

// 角色表
type Role struct {
	gorm.Model
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Users []User `gorm:"foreignKey:RoleID"`
}

// 用户表
type User struct {
	ID        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `json:"email"`
	Username  string `gorm:"unique" json:"username"`
	Picture   string `json:"picture"`
	Name      string `json:"name"`
	NickName  string `json:"nick_name"`
	Auth      Auth   `json:"auth"`
	RoleID    uint   `json:"role_id"`
	Role      Role   `gorm:"foreignKey:RoleID" json:"role"`
}

// 授权表
type Auth struct {
	ID           string `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	HashPassword string    `gorm:"not null" json:"-"` //- 表示不返回
	IsLogin      bool      `gorm:"default=false" json:"is_login"`
	Token        string    `gorm:"text" json:"token"`
	TokenExp     time.Time `gorm:"null" json:"token_exp"`
	UserID       string    `gorm:"unique" json:"user_id"`
	User         *User     `gorm:"foreignKey:UserID" json:"user"`
}

// 一级栏目表
type Category struct {
	ID            string `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string        `gorm:"not null;unique" json:"name"`
	SubCategories []SubCategory `gorm:"foreignKey:CategoryID;references:ID;" json:"sub_categories"`
}

// 二级分类
type SubCategory struct {
	ID         string `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string     `gorm:"not null;unique" json:"name"`
	CategoryID *string    `gorm:"default:null" json:"category_id"`
	Category   Category   `gorm:"foreignKey:CategoryID;" json:"category"`
	Properties []Property `gorm:"foreignKey:SubCategoryID;references:ID;" json:"properties"`
	Products   []Product  `gorm:"many2many:product_subcategories" json:"products"`
}

// 特性表
type Property struct {
	ID            string `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string      `json:"name"`
	SubCategoryID *string     `gorm:"default:null" json:"sub_category_id"`
	SubCategory   SubCategory `gorm:"foreignKey:SubCategoryID;" json:"sub_category"`
	Values        []Value     `gorm:"foreignKey:PropertyID;references:ID;" json:"values"`
	Products      []Product   `gorm:"many2many:product_properties" json:"products"`
}

// 特性值表
type Value struct {
	ID         string `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string   `json:"name"`
	PropertyID *string  `gorm:"default:null" json:"property_id"`
	Property   Property `gorm:"foreignKey:PropertyID;" json:"property"`
}

// 产品表
type Product struct {
	ID            string `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string         `gorm:"not null;unique" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         float64        `gorm:"not null" json:"price"`
	Stock         uint           `json:"stock"`
	Images        []*Image       `gorm:"many2many:product_images;" json:"images"`
	Properties    []*Property    `gorm:"many2many:product_properties;" json:"properties"`
	SubCategories []*SubCategory `gorm:"many2many:product_subcategories" json:"sub_categories"`
	IsShow        bool           `json:"is_show"`
	Start         time.Time      `json:"start"`
	End           time.Time      `json:"end"`
}

// 产品子栏目关系表
type ProductSubcategory struct {
	ProductID     string `gorm:"primaryKey"`
	SubCategoryID string `gorm:"primaryKey"`
}

// 产品图片关联表
type ProductImage struct {
	ProductID string `gorm:"primaryKey"`
	ImageID   string `gorm:"primaryKey"`
}

// 图片表
type Image struct {
	ID               string `gorm:"primaryKey" json:"uid"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Url              string    `json:"url"`
	Name             string    `json:"name"`
	Size             uint      `json:"size"`
	Status           string    `json:"status"`
	Type             string    `json:"type"`
	LastModifiedDate time.Time `gorm:"null;" json:"lastModifiedDate"`
}

// 产品图片关联表
type ProductProperty struct {
	ProductID  string `gorm:"primaryKey" json:"product_id"`
	PropertyID string `gorm:"primaryKey" json:"property_id"`
}

// 订单状态枚举
type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"   //创建
	OrderStatusPending   OrderStatus = "pending"   //待支付
	OrderStatusPaid      OrderStatus = "paid"      //已支付
	OrderStatusShipped   OrderStatus = "shipped"   //发货中
	OrderStatusCompleted OrderStatus = "completed" //已完成
	OrderStatusCancelled OrderStatus = "cancelled" //已取消
)

// 订单表
type Order struct {
	ID         string `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CustomerID string      `gorm:"not null" json:"customer_id"`                    // 关联客户表
	Total      float64     `gorm:"not null" json:"total"`                          //总金额
	Status     OrderStatus `gorm:"type:varchar(20);not null" json:"status"`        //订单状态
	CouponID   *string     `gorm:"default:null" json:"coupon_id"`                  // 关联优惠券表
	Items      []OrderItem `gorm:"foreignKey:OrderID;references:ID;" json:"items"` //订单项
}

// 订单项表
type OrderItem struct {
	ID        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderID   string  `gorm:"not null" json:"order_id"`   //订单ID
	ProductID string  `gorm:"not null" json:"product_id"` //产品ID
	Quantity  int     `gorm:"not null" json:"quantity"`   //购买数量
	Price     float64 `gorm:"not null" json:"price"`      //价格
}

// 购物车表
type Cart struct {
	ID         string `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CustomerID string     `gorm:"not null" json:"customer_id"`
	Items      []CartItem `gorm:"foreignKey:CartID;references:ID;" json:"items"` //购物车项
}

// 购物车项表
type CartItem struct {
	ID        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CartID    string `gorm:"not null" json:"cart_id"`    //购物车ID
	ProductID string `gorm:"not null" json:"product_id"` //产品ID
	Quantity  int    `gorm:"not null" json:"quantity"`   //产品数量
}

// 优惠券表
type Coupon struct {
	ID         string `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Code       string    `gorm:"not null;unique" json:"code"`
	Discount   float64   `gorm:"not null" json:"discount"`
	Expiration time.Time `gorm:"not null" json:"expiration"`
}

// 客户表
type Customer struct {
	ID        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string   `gorm:"not null" json:"first_name"`
	LastName  string   `gorm:"not null" json:"last_name"`
	Email     string   `gorm:"not null;unique" json:"email"`
	Phone     string   `gorm:"not null" json:"phone"`
	Address   string   `gorm:"not null" json:"address"`
	Coupons   []Coupon `gorm:"many2many:customer_coupons;" json:"coupons"`
}

// 客户优惠券关联表
type CustomerCoupon struct {
	CustomerID string `gorm:"primaryKey"`
	CouponID   string `gorm:"primaryKey"`
}
