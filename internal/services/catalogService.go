package services

import (
	"errors"
	"go-ecommerce-app/configs"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config configs.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequestDto) error {
	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageURL,
		DisplayOrder: input.DisplayOrder,
	})

	return err
}

func (s CatalogService) EditCategory(id int, input dto.CreateCategoryRequestDto) (*domain.Category, error) {
	existCat, err := s.Repo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if len(input.Name) > 0 {
		existCat.Name = input.Name
	}
	if input.ParentId > 0 {
		existCat.ParentId = input.ParentId
	}
	if len(input.ImageURL) > 0 {
		existCat.ImageUrl = input.ImageURL
	}
	if input.DisplayOrder > 0 {
		existCat.DisplayOrder = input.DisplayOrder
	}

	updatedCat, err := s.Repo.EditCategory(existCat)

	return updatedCat, err
}

func (s CatalogService) DeleteCategory(id int) error {
	err := s.Repo.DeleteCategory(id)
	if err != nil {
		return errors.New("category not found to delete")
	}

	return nil
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, errors.New("no categories found")
	}

	return categories, nil
}

func (s CatalogService) GetCategory(id int) (*domain.Category, error) {
	cat, err := s.Repo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	return cat, nil
}

////// Products ///////

func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryId:  input.CategoryId,
		ImageUrl:    input.ImageUrl,
		UserId:      int(user.ID),
		Stock:       uint(input.Stock),
	})

	return err
}

func (s CatalogService) EditProduct(id int, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {
	existProduct, err := s.Repo.FindProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if existProduct.UserId != int(user.ID) {
		return nil, errors.New("you are not authorized to update this product")
	}

	if len(input.Name) > 0 {
		existProduct.Name = input.Name
	}
	if len(input.Description) > 0 {
		existProduct.Description = input.Description
	}
	if input.CategoryId > 0 {
		existProduct.CategoryId = input.CategoryId
	}
	if input.Price > 0 {
		existProduct.Price = input.Price
	}

	updatedProduct, err := s.Repo.EditProduct(existProduct)

	return updatedProduct, err
}

func (s CatalogService) DeleteProduct(id int, user domain.User) error {
	existProduct, err := s.Repo.FindProductByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	if existProduct.UserId != int(user.ID) {
		return errors.New("you are not authorized to delete this product")
	}

	err = s.Repo.DeleteProduct(int(existProduct.ID))
	if err != nil {
		return errors.New("product not found to delete")
	}

	return nil
}

func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, errors.New("no products found")
	}

	return products, nil
}

func (s CatalogService) GetProductById(id int) (*domain.Product, error) {
	product, err := s.Repo.FindProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return product, nil
}

func (s CatalogService) GetSellerProducts(id int) ([]*domain.Product, error) {
	products, err := s.Repo.FindSellerProducts(id)
	if err != nil {
		return nil, errors.New("no products found for this seller")
	}

	return products, err
}

func (s CatalogService) UpdateProductStock(e domain.Product) (*domain.Product, error) {
	product, err := s.Repo.FindProductByID(int(e.ID))
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.UserId != e.UserId {
		return nil, errors.New("you are not authorized to update this product")
	}

	product.Stock = e.Stock
	editProduct, err := s.Repo.EditProduct(product)
	if err != nil {
		return nil, err
	}

	return editProduct, err
}
