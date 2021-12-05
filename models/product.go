package models

import "github.com/google/uuid"

type Product struct {
	BaseModel
	Name  string  `json:"name"`
	Price float64 `gorm:"price" json:"price"`
}

type Order struct {
	BaseModel
	UserID uuid.UUID `gorm:"column:user_id;type:char(36)" json:"user_id"`
	Total  float64   `gorm:"total" json:"total"`
}

type OrderItem struct {
	BaseModel
	ProductID uuid.UUID `gorm:"column:product_id;type:char(36)" json:"product_id"`
	Quantity  int64     `gorm:"qty" json:"quantity"`
	Price     float64   `gorm:"price" json:"price"`
	Disc      float64   `gorm:"disc" json:"disc"`
	Total     float64   `gorm:"total_price" json:"total_price"`
}

type Stock struct {
	BaseModel
	ProductID uuid.UUID `gorm:"column:product_id;type:char(36)" json:"product_id"`
	Quantity  int64     `gorm:"qty" json:"qty"`
	Remark    string    `json:"remark"`
}

// this will be database for reserved item
type ReservedStock struct {
	BaseModel
	ProductID uuid.UUID `gorm:"column:product_id;type:char(36)" json:"product_id"`
	CartID    uuid.UUID `gorm:"column:cart_id;type:char(36)" json:"cart_id"`
	Quantity  int64     `gorm:"qty" json:"qty"`
	Remark    string    `json:"remark"`
}
