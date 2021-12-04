package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB *gorm.DB
}

func (h CartHandler) AddNewItemToCart(c *fiber.Ctx) error {

	return nil
}
