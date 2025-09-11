package handlers

import (
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/services"
	"strconv"

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
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProduct)
	app.Get("/categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategoryById)

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
func (h *CatalogHandler) GetCategories(ctx *fiber.Ctx) error {
	cats, err := h.svc.GetCategories()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "GetCategories", cats)
}

func (h *CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	cat, err := h.svc.GetCategory(id)
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "GetCategoryById", cat)
}

func (h *CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequestDto{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "create category request body is not valid")
	}

	err = h.svc.CreateCategory(req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessMessage(ctx, "category created succesfully", nil)
}

func (h *CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	req := dto.CreateCategoryRequestDto{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "update category request body is not valid")
	}

	updatedCat, err := h.svc.EditCategory(id, req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessMessage(ctx, "category updated succesfully", updatedCat)
}

func (h *CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := h.svc.DeleteCategory(id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessMessage(ctx, "category deleted successfully", nil)
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
