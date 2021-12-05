package services

import (
	"fmt"
	"log"
	"time"

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

// this will run every time set on cron at the main.go, and reset the stock back to stock table
func (s ProductService) ResetStock(t time.Time) {
	var reservedStocks []models.ReservedStock

	s.DB.Find(&reservedStocks)
	log.Printf("reservedStocks %v", reservedStocks)
	for _, reservedStock := range reservedStocks {
		stock := models.Stock{
			ProductID: reservedStock.ProductID,
		}
		_ = s.GetStockByProductId(&stock)
		stock.Quantity = stock.Quantity + reservedStock.Quantity
		s.DB.Updates(&stock)
		s.DB.Delete(&reservedStock)
	}
}
