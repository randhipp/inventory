package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WebhookHandler struct {
	DB *gorm.DB
}

func (h WebhookHandler) NewPaymentWebhook(c *fiber.Ctx) error {

	return nil
}
