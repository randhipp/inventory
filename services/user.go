package services

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/randhipp/inventory/models"
)

type UserService struct {
	DB *gorm.DB
}

func (s UserService) GetUserByID(user *models.User) error {
	if err := s.DB.Preload("Merchant").First(user).Error; err != nil {
		fmt.Println("GetUserByID Err")
		fmt.Println(err)
		return err
	}
	return nil
}

func (s UserService) GetUserByEmail(user *models.User) error {
	if err := s.DB.Preload("Merchant").First(user, "email = ?", user.Email).Error; err != nil {
		fmt.Println("GetUserByID Err")
		fmt.Println(err)
		return err
	}
	return nil
}
