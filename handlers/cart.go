package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/randhipp/inventory/models"
	"github.com/randhipp/inventory/services"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB      *gorm.DB
	Product services.ProductService
}

func (h CartHandler) AddNewItemToCart(c *fiber.Ctx) error {
	req := &models.CartRequest{}
	if err := c.BodyParser(req); err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(models.Error{
			Message: "invalid payload",
			Field:   "*",
		})
		return nil
	}

	// get user from jwt token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	userId, _ := uuid.Parse(claims["id"].(string))

	cart := models.Cart{
		UserID: userId,
	}

	total := 0.0
	cartItem := req.CartItem
	product := models.Product{}
	product.ID = cartItem.ProductID

	err := h.Product.GetProductById(&product)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.ErrUnprocessableEntity.Code).JSON(models.Error{
			Message: "product unavailable",
		})
	}

	stock := models.Stock{
		ProductID: cartItem.ProductID,
	}
	err = h.Product.GetStockByProductId(&stock)
	if err != nil {
		fmt.Println(err)
	}
	if stock.Quantity < 1 || stock.Quantity < cartItem.Quantity {
		return c.Status(fiber.ErrUnprocessableEntity.Code).JSON(models.Error{
			Message: "stock unavailable",
		})
	}

	cart.Total = total
	h.DB.Transaction(func(tx *gorm.DB) error {
		// create cart
		if err := tx.Create(&cart).Error; err != nil {
			// return any error will rollback
			fmt.Println("failed new cart transaction")
			fmt.Println(err)
			c.Status(fiber.ErrInternalServerError.Code).JSON(models.Error{
				Message: fiber.ErrInternalServerError.Message,
			})
			return err
		}

		// create new item
		newItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
			Disc:      0,
			Total:     product.Price * float64(cartItem.Quantity),
		}
		if err := tx.Create(&newItem).Error; err != nil {
			// return any error will rollback
			fmt.Println("failed new cart item transaction")
			fmt.Println(err)
			c.Status(fiber.ErrInternalServerError.Code).JSON(models.Error{
				Message: fiber.ErrInternalServerError.Message,
			})
			return err
		}

		// reduce stock
		fmt.Println(stock)
		stock.Quantity = stock.Quantity - cartItem.Quantity
		fmt.Println(stock)
		tx.Save(&stock)

		// add stock to reserved table, we can revert back using scheduler for each amount of time
		reservedStock := models.ReservedStock{
			ProductID: product.ID,
			CartID:    cart.ID,
			Quantity:  cartItem.Quantity,
			Remark:    fmt.Sprintf("reserved for user : %s", name),
		}
		if err := tx.Create(&reservedStock).Error; err != nil {
			// return any error will rollback
			fmt.Println("failed reservedStock transaction")
			fmt.Println(err)
			c.Status(fiber.ErrInternalServerError.Code).JSON(models.Error{
				Message: fiber.ErrInternalServerError.Message,
			})
			return err
		}
		return nil
	})
	c.JSON(models.CartResponse{
		Status: "success",
		Cart:   cart,
	})
	return nil
}
