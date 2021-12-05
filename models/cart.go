package models

import "github.com/google/uuid"

type Cart struct {
	BaseModel
	UserID   uuid.UUID  `gorm:"column:user_id;type:char(36)" json:"user_id"`
	Total    float64    `gorm:"total" json:"total"`
	CartItem []CartItem `json:"cart_items,omitempty"`
}

type CartItem struct {
	BaseModel
	CartID    uuid.UUID `gorm:"column:cart_id;type:char(36)" json:"cart_id"`
	ProductID uuid.UUID `gorm:"column:product_id;type:char(36)" json:"product_id"`
	Quantity  int64     `gorm:"qty" json:"qty"`
	Price     float64   `gorm:"price" json:"price"`
	Disc      float64   `gorm:"disc" json:"disc"`
	Total     float64   `gorm:"total_price" json:"total_price"`
}

type CartRequest struct {
	CartItem
}

type CartResponse struct {
	Status string `json:"status"`
	Cart   Cart   `json:"cart"`
}
