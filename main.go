package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
	"github.com/randhipp/inventory/database"
	"github.com/randhipp/inventory/handlers"
	"github.com/randhipp/inventory/models"
	"github.com/randhipp/inventory/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func jwtCustomError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.ErrUnauthorized.Code).JSON(models.LoginResponse{
		Status: fiber.ErrUnauthorized.Message,
	})
}

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	AuthHandler := handlers.AuthHandler{
		DB: database.DBConn,
		AuthService: services.AuthService{
			DB: database.DBConn,
		},
		UserService: services.UserService{
			DB: database.DBConn,
		},
	}
	app.Post("/api/v1/auth", AuthHandler.Login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("HMAC_SECRET")),
		ErrorHandler: jwtCustomError,
	}))

	UserHandler := handlers.UserHandler{
		DB: database.DBConn,
		UserService: services.UserService{
			DB: database.DBConn,
		},
	}
	app.Get("/api/v1/users", UserHandler.GetUsers)
	app.Get("/api/v1/users/:id", UserHandler.GetUser)
	app.Post("/api/v1/users", UserHandler.NewUser)
	app.Patch("/api/v1/users/:id", UserHandler.UpdateUser)
	app.Delete("/api/v1/users/:id", UserHandler.DeleteUser)

	CartHandler := handlers.CartHandler{
		DB: database.DBConn,
	}
	PaymentHandler := handlers.PaymentHandler{
		DB: database.DBConn,
	}
	WebhookHandler := handlers.WebhookHandler{
		DB: database.DBConn,
	}

	// we will simulate user add to cart and payment using this 3 API
	app.Post("/api/v1/cart", CartHandler.AddNewItemToCart)
	app.Post("/api/v1/cart/:cartId/pay", PaymentHandler.NewPayment)
	app.Post("/api/v1/webhooks/payments/:paymentId", WebhookHandler.NewPaymentWebhook)
}

func initDatabase() {
	var err error
	dsn := fmt.Sprintf("%v:%v@%v/%v?multiStatements=true&parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_NAME"),
	)
	database.DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection Opened to Database")
	err = database.DBConn.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err)
		panic("failed to migrate User")
	}
	err = database.DBConn.AutoMigrate(&models.Product{})
	if err != nil {
		fmt.Println(err)
		panic("failed to migrate User")
	}
	err = database.DBConn.AutoMigrate(&models.Order{})
	if err != nil {
		fmt.Println(err)
		panic("failed to migrate User")
	}
	err = database.DBConn.AutoMigrate(&models.OrderItem{})
	if err != nil {
		fmt.Println(err)
		panic("failed to migrate User")
	}
	err = database.DBConn.AutoMigrate(&models.Stock{})
	if err != nil {
		fmt.Println(err)
		panic("failed to migrate User")
	}
	err = database.DBConn.AutoMigrate(&models.StockReserved{})
	if err != nil {
		fmt.Println(err)
		panic("failed to migrate User")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	var users = []models.User{
		{
			Name:     "admin1",
			Email:    "admin1@admin.com",
			Password: string(hashedPassword),
		},
		{
			Name:     "user1",
			Email:    "user1@user.com",
			Password: string(hashedPassword),
		},
		{
			Name:     "user2",
			Email:    "user2@user.com",
			Password: string(hashedPassword),
		},
		{
			Name:     "user3",
			Email:    "user3@user.com",
			Password: string(hashedPassword),
		},
		{
			Name:     "user4",
			Email:    "user4@user.com",
			Password: string(hashedPassword),
		},
	}
	database.DBConn.Unscoped().Delete(&models.User{}, "email LIKE ?", "%")
	database.DBConn.Create(&users)

	// seed for product
	productDemo := models.Product{
		Name:  "Apple iPhone 13 Pro Max",
		Price: 25000000,
	}
	productDemo.ID, err = uuid.Parse("3829cf54-2680-4728-a53f-7cea064f12be")
	if err != nil {
		fmt.Println(err)
		panic("failed to parse UUID")
	}
	var products = []models.Product{productDemo}

	database.DBConn.Unscoped().Delete(&models.Product{}, "name LIKE ?", "%")
	database.DBConn.Create(&products)

	// seed for stock
	stock := models.Stock{
		ProductID: productDemo.ID,
		Quantity:  1,
	}
	var stocks = []models.Stock{stock}
	database.DBConn.Unscoped().Delete(&models.Stock{}, "id LIKE ?", "%")
	database.DBConn.Create(&stocks)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	app.Listen(":3000")
}
