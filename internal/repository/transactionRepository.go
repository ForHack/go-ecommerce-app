package repository

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

type transactionStorage struct {
	db *gorm.DB
}

func (t *transactionStorage) CreatePayment(payment *domain.Payment) error {
	return nil
}

func (t *transactionStorage) FindOrders(uId uint) ([]domain.OrderItem, error) {
	return nil, nil
}

func (t *transactionStorage) FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error) {
	return dto.SellerOrderDetails{}, nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionStorage{db: db}
}
