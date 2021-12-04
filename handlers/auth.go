package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/randhipp/inventory/models"
	"github.com/randhipp/inventory/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB          *gorm.DB
	AuthService services.AuthService
	UserService services.UserService
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	loginRequest := &models.LoginRequest{}
	if err := c.BodyParser(loginRequest); err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(models.Error{
			Message: "invalid payload",
			Field:   "*",
		})
		return nil
	}
	user := models.User{
		Email: loginRequest.Email,
	}
	err := h.UserService.GetUserByEmail(&user)
	if err != nil {
		fmt.Println(err)
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.ErrUnauthorized.Code).JSON(models.LoginResponse{
			Status: fiber.ErrUnauthorized.Message,
		})
	}
	tokenString, err := h.AuthService.GetToken(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c.JSON(models.LoginResponse{
		Status: "success",
		Token:  tokenString,
	})
	return nil
}
