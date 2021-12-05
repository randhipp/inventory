package services

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/randhipp/inventory/models"
)

type ProductService struct {
	DB *gorm.DB
}

func (s ProductService) GetProductById(product *models.Product) error {
	if err := s.DB.First(product).Error; err != nil {
		fmt.Println("GetProductById Err")
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ProductService) GetStockByProductId(stock *models.Stock) error {
	if err := s.DB.First(stock, "product_id = ?", stock.ProductID).Error; err != nil {
		fmt.Println("GetStockByProductId Err")
		fmt.Println(err)
		return err
	}
	return nil
}
