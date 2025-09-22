package handlers

import (
	"errors"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/services"
	"go-ecommerce-app/pkg/payment"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	svc           services.TransactionService
	userSvc       services.UserService
	paymentClient payment.PaymentClient
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth) services.TransactionService {
	return services.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(as *rest.RestHandler) {
	app := as.App
	svc := initializeTransactionService(as.DB, as.Auth)
	userSvc := services.UserService{
		Repo:   repository.NewUserRepository(as.DB),
		CRepo:  repository.NewCatalogRepository(as.DB),
		Auth:   as.Auth,
		Config: as.Config,
	}

	handler := &TransactionHandler{
		svc:           svc,
		paymentClient: as.Pc,
		userSvc:       userSvc,
	}

	secRoute := app.Group("/", as.Auth.Authorize)
	secRoute.Get("/payment", handler.MakePayment)

	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {
	// gram authorize user
	user := h.svc.Auth.GetCurrentUser(ctx)

	activePayment, err := h.svc.GetActivePayment(user.ID)
	if activePayment.ID > 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"error":       "You have an ongoing payment. Please complete it before initiating a new one.",
			"payment_url": activePayment.PaymentUrl,
		})
	}

	_, amount, err := h.userSvc.FindCart(user.ID)

	orderId, err := helper.RandomNumbers(8)
	if err != nil {
		return rest.InternalError(ctx, errors.New("failed to generate order id"))
	}

	sessionResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)

	err = h.svc.StoreCreatedPayment(user.ID, sessionResult, amount, orderId)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Payment initiated successfully",
		"result":      sessionResult,
		"payment_url": sessionResult.URL,
	})
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("Get all orders")
}

func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("GetOrderDetails")
}
