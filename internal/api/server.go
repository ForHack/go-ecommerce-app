package api

import (
	"go-ecommerce-app/configs"

	"github.com/gofiber/fiber/v2"
)

func StartServer(config configs.AppConfig) {
	app := fiber.New()

	app.Get("/health", HealthCheck)

	app.Listen(config.ServerPort)
}

func HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "API is up and running",
	})
}
