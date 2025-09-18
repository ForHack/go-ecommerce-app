package api

import (
	"go-ecommerce-app/configs"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config configs.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// run migration
	err = db.
		AutoMigrate(
			&domain.User{},
			&domain.Address{},
			&domain.BankAccount{},
			&domain.Category{},
			&domain.Product{},
			&domain.Cart{},
			&domain.Order{},
			&domain.OrderItem{},
			&domain.Payment{},
		)
	if err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	c := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	})
	app.Use(c)

	auth := helper.SetupAuth(config.AppSecret)

	rh := &rest.RestHandler{App: app, DB: db, Auth: auth, Config: config}
	setupRoutes(rh)

	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)
	handlers.SetupTransactionRoutes(rh)
	handlers.SetupCatalogRoutes(rh)
}
