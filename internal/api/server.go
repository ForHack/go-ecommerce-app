package api

import (
	"go-ecommerce-app/configs"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/handlers"

	"github.com/gofiber/fiber/v2"
)

func StartServer(config configs.AppConfig) {
	app := fiber.New()

	restHandler := &rest.RestHandler{App: app}
	setupRoutes(restHandler)

	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.RestHandler) {
	// user handler
	handlers.SetupUserRoutes(rh)
}
