package handlers

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/randhipp/inventory/helpers"
	"github.com/randhipp/inventory/models"
	"github.com/randhipp/inventory/services"
	"github.com/vcraescu/go-paginator/v2"
	"github.com/vcraescu/go-paginator/v2/adapter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB          *gorm.DB
	UserService services.UserService
}

func (h UserHandler) GetUsers(c *fiber.Ctx) error {
	var users []models.User
	q := h.DB.Preload("Merchant").Model(models.User{})

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	p := paginator.New(adapter.NewGORMAdapter(q), limit)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	p.SetPage(page)

	if err := p.Results(&users); err != nil {
		panic(err)
	}

	pages, _ := p.PageNums()
	current, _ := p.Page()

	resp := models.UsersResponse{
		Status:      "success",
		Users:       users,
		TotalPage:   pages,
		CurrentPage: current,
		PerPage:     limit,
	}
	next, err := p.NextPage()
	if err == nil {
		resp.NextPage = next
	}
	prev, err := p.PrevPage()
	if err == nil {
		resp.PrevPage = prev
	}
	c.JSON(resp)
	return nil
}

func (h UserHandler) GetUser(c *fiber.Ctx) error {
	userId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	user := &models.User{}
	user.ID = userId
	err = h.UserService.GetUserByID(user)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.ErrNotFound.Code).JSON(models.Error{
			Message: fiber.ErrNotFound.Message,
			Field:   "ID",
		})
	}
	resp := models.UserResponse{
		Status: "success",
		User:   user,
	}
	c.JSON(resp)
	return nil
}

func (h UserHandler) NewUser(c *fiber.Ctx) error {
	userRequest := &models.UserRequest{}
	u := &models.User{}
	if err := c.BodyParser(userRequest); err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(models.Error{
			Message: "invalid payload",
			Field:   "*",
		})
		return nil
	}

	if helpers.IsEmailValid(userRequest.Email) == false {
		c.Status(fiber.ErrConflict.Code).JSON(models.Error{
			Message: "invalid email",
			Field:   "email",
		})
		return nil
	}

	if len(userRequest.Password) < 8 {
		c.Status(fiber.ErrConflict.Code).JSON(models.Error{
			Message: "password require 8 or more character",
			Field:   "password",
		})
		return nil
	}

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	u.Email = userRequest.Email
	u.Name = userRequest.Name
	u.Password = string(hashedPassword)

	h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&u).Error; err != nil {
			// return any error will rollback
			fmt.Println("failed new user transaction")
			fmt.Println(err)
			c.Status(fiber.ErrConflict.Code).JSON(models.Error{
				Message: "conflict",
				Field:   "email",
			})
			return err
		}
		c.JSON(models.UserResponse{
			Status: "success",
			User:   u,
		})
		return nil
	})
	return nil
}

func (h UserHandler) UpdateUser(c *fiber.Ctx) error {
	userId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	u := models.User{}
	u.ID = userId
	err = h.UserService.GetUserByID(&u)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.ErrNotFound.Code).JSON(models.Error{
			Message: fiber.ErrNotFound.Message,
			Field:   "ID",
		})
	}
	userRequest := &models.UserRequest{}
	if err := c.BodyParser(userRequest); err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(models.Error{
			Message: "invalid payload",
			Field:   "*",
		})
		return nil
	}

	if userRequest.Email != "" {
		if helpers.IsEmailValid(userRequest.Email) == false {
			c.Status(fiber.ErrConflict.Code).JSON(models.Error{
				Message: "invalid email",
				Field:   "email",
			})
			return nil
		}
		u.Email = userRequest.Email
	}

	if userRequest.Name != "" {
		u.Name = userRequest.Name
	}

	if userRequest.Password != "" {
		if len(userRequest.Password) < 8 {
			c.Status(fiber.ErrConflict.Code).JSON(models.Error{
				Message: "password require 8 or more character",
				Field:   "password",
			})
			return nil
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
		u.Password = string(hashedPassword)
	}

	h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(&u).Error; err != nil {
			// return any error will rollback
			fmt.Println("failed new user transaction")
			fmt.Println(err)
			c.Status(fiber.ErrInternalServerError.Code).JSON(models.Error{
				Message: fiber.ErrInternalServerError.Message,
			})
			return err
		}
		c.JSON(models.UserResponse{
			Status: "success",
			User:   &u,
		})
		return nil
	})
	return nil
}

func (h UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	h.DB.First(&user, "id = ?", id)
	if user.Email == "" {
		c.Status(404).JSON(models.UserDeleteResponse{
			Status: "Not Found",
		})
		return nil
	}
	h.DB.Delete(&user)
	c.Status(202).JSON(models.UserDeleteResponse{
		Status: "User Successfully deleted",
	})
	return nil
}
