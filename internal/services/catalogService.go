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
