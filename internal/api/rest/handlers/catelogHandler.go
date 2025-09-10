package handlers

import (
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

type CatalogHandler struct {
	svc services.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	svc := services.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}

	handler := &CatalogHandler{
		svc: svc,
	}

	// Public routes
	app.Get("/products")
	app.Get("/products/:id")
	app.Get("/categories")
	app.Get("/categories/:id")

	// Protected routes
	sellerRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)

	// Category routes
	sellerRoutes.Post("/categories", handler.CreateCategories)
	sellerRoutes.Patch("/categories/:id", handler.EditCategory)
	sellerRoutes.Delete("/categories/:id", handler.DeleteCategory)

	// Product routes
	sellerRoutes.Post("/products", handler.CreateProducts)
	sellerRoutes.Get("/products", handler.GetProducts)
	sellerRoutes.Get("/products/:id", handler.GetProduct)
	sellerRoutes.Put("/products/:id", handler.EditProduct)
	sellerRoutes.Patch("/products/:id", handler.UpdateStock)
	sellerRoutes.Delete("/products/:id", handler.DeleteProduct)
}

// /////////////////////////// Categories /////////////////////////////////////
func (h *CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	log.Println("Current User: ", user)

	return rest.SuccessMessage(ctx, "CreateCategories", nil)
}

func (h *CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "EditCategory", nil)
}

func (h *CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "DeleteCategory", nil)
}

///////////////////////////// Products /////////////////////////////////////

func (h *CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "CreateProducts", nil)
}

func (h *CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "CreateProducts", nil)
}

func (h *CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "GetProducts", nil)
}

func (h *CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "GetProduct", nil)
}

func (h *CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "UpdateStock", nil)
}

func (h *CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	return rest.SuccessMessage(ctx, "DeleteProduct", nil)
}
