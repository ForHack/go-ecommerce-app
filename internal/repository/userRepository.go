package repository

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(u domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserByID(id uint) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)

	CreateBankAccount(e domain.BankAccount) error

	FindCartItems(uId uint) ([]domain.Cart, error)
	FindCartItem(uId uint, pId uint) (domain.Cart, error)
	CreateCart(e domain.Cart) error
	UpdateCart(e domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(uId uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r userRepository) CreateUser(user domain.User) (domain.User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return domain.User{}, errors.New("failed to create user")
	}

	return user, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User

	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		log.Printf("Find user error %v", err)
		return domain.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (r userRepository) FindUserByID(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.First(&user, id).Error
	if err != nil {
		log.Printf("Find user error %v", err)
		return domain.User{}, errors.New("user does not exist")
	}

	return user, nil
}

func (r userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {
	var user domain.User

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(u).Error
	if err != nil {
		log.Printf("Update user error %v", err)
		return domain.User{}, errors.New("failed to update user")
	}

	return user, nil
}

func (r userRepository) CreateBankAccount(e domain.BankAccount) error {
	return r.db.Create(&e).Error
}

// CreateCart implements UserRepository.
func (r *userRepository) CreateCart(c domain.Cart) error {
	return r.db.Create(&c).Error
}

// DeleteCartById implements UserRepository.
func (r *userRepository) DeleteCartById(id uint) error {
	err := r.db.Delete(&domain.Cart{}, id).Error
	return err
}

// DeleteCartItems implements UserRepository.
func (r *userRepository) DeleteCartItems(uId uint) error {
	err := r.db.Where("user_id = ?", uId).Delete(&domain.Cart{}).Error
	return err
}

// FindCartItem implements UserRepository.
func (r *userRepository) FindCartItem(uId uint, pId uint) (domain.Cart, error) {
	cartItem := domain.Cart{}
	err := r.db.Where("user_id = ? AND product_id=?", uId, pId).First(&cartItem).Error
	return cartItem, err
}

// FindCartItems implements UserRepository.
func (r *userRepository) FindCartItems(uId uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Where("user_id = ?", uId).Find(&carts).Error
	return carts, err
}

// UpdateCart implements UserRepository.
func (r *userRepository) UpdateCart(c domain.Cart) error {
	var cart domain.Cart
	err := r.db.Model(&cart).Clauses(clause.Returning{}).Where("id = ?", c.ID).Updates(c).Error
	return err
}
