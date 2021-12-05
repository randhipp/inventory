package models

import "github.com/google/uuid"

type Product struct {
	BaseModel
	Name  string  `json:"name"`
	Price float64 `gorm:"price" json:"price"`
}

type Order struct {
	BaseModel
	Total float64 `gorm:"total" json:"total"`
}

type OrderItem struct {
	BaseModel
	ProductID uuid.UUID `gorm:"column:product_id;type:char(36)" json:"productId"`
	Quantity  int64     `gorm:"qty" json:"quantity"`
	Price     float64   `gorm:"price" json:"price"`
	Disc      float64   `gorm:"disc" json:"disc"`
	Total     float64   `gorm:"totalPrice" json:"totalPrice"`
}

type Stock struct {
	BaseModel
	ProductID uuid.UUID `gorm:"column:product_id;type:char(36)" json:"productId"`
	Quantity  int64     `gorm:"qty" json:"quantity"`
	Remark    string    `json:"remark"`
}

// this will be database for reserved item
type StockReserved struct {
	BaseModel
	UserID    uuid.UUID `gorm:"column:user_id;type:char(36)" json:"user_id"`
	ProductID uuid.UUID `gorm:"column:product_id;type:char(36)" json:"productId"`
	OrderId   uuid.UUID `gorm:"column:order_id;type:char(36)" json:"orderId"`
	Quantity  int64     `gorm:"qty" json:"quantity"`
	Remark    string    `json:"remark"`
}
