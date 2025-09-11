package repository

import (
	"errors"
	"go-ecommerce-app/internal/domain"

	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateCategory(e *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryByID(id int) (*domain.Category, error)
	EditCategory(e *domain.Category) (*domain.Category, error)
	DeleteCategory(id int) error

	// Product methods can be added here
	CreateProduct(e *domain.Product) error
	FindProducts() ([]*domain.Product, error)
	FindProductByID(id int) (*domain.Product, error)
	FindSellerProducts(id int) ([]*domain.Product, error)
	EditProduct(e *domain.Product) (*domain.Product, error)
	DeleteProduct(id int) error
}

type catalogRepository struct {
	db *gorm.DB
}

func (c *catalogRepository) CreateCategory(e *domain.Category) error {
	err := c.db.Create(e).Error
	if err != nil {
		return errors.New("failed to create category: ")
	}
	return nil
}

func (c *catalogRepository) FindCategories() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := c.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *catalogRepository) FindCategoryByID(id int) (*domain.Category, error) {
	var category *domain.Category

	err := c.db.First(&category, id).Error
	if err != nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

func (r *catalogRepository) EditCategory(e *domain.Category) (*domain.Category, error) {
	err := r.db.Save(e).Error

	if err != nil {
		return nil, errors.New("failed to update category")
	}

	return e, nil
}

func (r *catalogRepository) DeleteCategory(id int) error {
	err := r.db.Delete(&domain.Category{}, id).Error

	if err != nil {
		return errors.New("failed to delete category")
	}

	return nil
}

// ///////////////////////////// Products /////////////////////////////////////
func (c *catalogRepository) CreateProduct(e *domain.Product) error {
	err := c.db.Model(&domain.Product{}).Create(e).Error
	if err != nil {
		return errors.New("failed to create product: ")
	}

	return nil
}

func (c *catalogRepository) FindProducts() ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c *catalogRepository) FindProductByID(id int) (*domain.Product, error) {
	var product *domain.Product
	err := c.db.First(&product, id).Error
	if err != nil {
		return nil, errors.New("product not found")
	}

	return product, nil
}

func (c *catalogRepository) FindSellerProducts(id int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Where("user_id = ?", id).Find(&products).Error
	if err != nil {
		return nil, errors.New("no products found for this seller")
	}

	return nil, nil
}

func (r *catalogRepository) EditProduct(e *domain.Product) (*domain.Product, error) {
	err := r.db.Save(e).Error
	if err != nil {
		return nil, errors.New("failed to update product")
	}

	return e, nil
}

func (r *catalogRepository) DeleteProduct(id int) error {
	err := r.db.Delete(&domain.Product{}, id).Error
	if err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{db: db}
}
