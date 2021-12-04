package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB *gorm.DB
}

func (h PaymentHandler) NewPayment(c *fiber.Ctx) error {

	return nil
}
