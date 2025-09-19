package rest

import (
	"go-ecommerce-app/configs"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/pkg/payment"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config configs.AppConfig
	Pc     payment.PaymentClient
}
