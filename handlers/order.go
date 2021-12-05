package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

func (h OrderHandler) NewOrder(c *fiber.Ctx) error {

	return nil
}
