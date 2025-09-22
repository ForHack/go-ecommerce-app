package handlers

import (
	"errors"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc services.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App

	svc := services.UserService{
		Repo:   repository.NewUserRepository(rh.DB),
		CRepo:  repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}
	handler := &UserHandler{
		svc: svc,
	}

	pubRoutes := app.Group("/users")
	// Public routes
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)

	pvtRoutes := pubRoutes.Group("/", rh.Auth.Authorize)

	// Protected routes
	pvtRoutes.Get("/verify", handler.GetVerificationCode)
	pvtRoutes.Post("/verify", handler.Verify)

	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Get("/profile", handler.GetProfile)
	pvtRoutes.Patch("/profile", handler.UpdateProfile)

	pvtRoutes.Post("/cart", handler.AddToCart)
	pvtRoutes.Get("/cart", handler.GetCart)

	pvtRoutes.Post("/order", handler.CreateOrder)
	pvtRoutes.Get("/order", handler.GetOrders)
	pvtRoutes.Get("/order/:id", handler.GetOrder)

	pvtRoutes.Post("/become-seller", handler.BecomeSeller)
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignUp{}
	err := ctx.BodyParser(&user)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid details",
		})
	}

	token, err := h.svc.SignUp(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error on creating account",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Account created successfully",
		"token":   token,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	loginInput := dto.UserLogin{}
	err := ctx.BodyParser(&loginInput)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid details",
		})
	}

	token, err := h.svc.Login(loginInput.Email, loginInput.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid credentials",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Login",
		"token":   token,
	})
}

func (h *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	// Create verification code and update to user profile in DB
	err := h.svc.GetVerificationCode(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Unable to generate verification code",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "GetVerificationCode",
	})
}
func (h *UserHandler) Verify(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	// Request
	var req dto.VerificationCodeInput

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid details",
		})
	}

	err := h.svc.VerifyCode(user.ID, req.Code)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Verified Successfully",
	})
}

func (h *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	err = h.svc.CreateProfile(user.ID, req)

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "CreateProfile",
	})
}
func (h *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "GetProfile",
		"profile": profile,
	})
}
func (h *UserHandler) UpdateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	err = h.svc.UpdateProfile(user.ID, req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "profile updated successfully",
	})
}

func (h *UserHandler) AddToCart(ctx *fiber.Ctx) error {
	req := dto.CreateCartRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	user := h.svc.Auth.GetCurrentUser(ctx)

	// Call service to add to cart
	cartItems, err := h.svc.CreateCart(req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessMessage(ctx, "Product added to cart successfully", cartItems)
}
func (h *UserHandler) GetCart(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	cart, _, err := h.svc.FindCart(user.ID)

	if err != nil {
		return rest.InternalError(ctx, errors.New("cart does not exist"))
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "GetCart",
		"cart":    cart,
	})
}

func (h *UserHandler) CreateOrder(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	orderRef, err := h.svc.CreateOrder(user)
	if err != nil {
		return rest.InternalError(ctx, errors.New("failed to create order"))
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "order created successfully",
		"order":   orderRef,
	})
}
func (h *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	orders, err := h.svc.GetOrders(user)
	if err != nil {
		return rest.InternalError(ctx, errors.New("failed to fetch orders"))
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "get orders",
		"orders":  orders,
	})
}
func (h *UserHandler) GetOrder(ctx *fiber.Ctx) error {
	orderId, _ := strconv.Atoi(ctx.Params("id"))
	user := h.svc.Auth.GetCurrentUser(ctx)

	order, err := h.svc.GetOrderById(uint(orderId), user.ID)
	if err != nil {
		return rest.InternalError(ctx, errors.New("failed to fetch order"))
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "GetOrder",
		"order":   order,
	})
}

func (h *UserHandler) BecomeSeller(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.SellerInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	token, err := h.svc.BecomeSeller(user.ID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "failed to become a seller",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Become seller successfully",
		"token":   token,
	})
}
