package services

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

type UserService struct {
	Repo repository.UserRepository
	Auth helper.Auth
}

func (s *UserService) SignUp(input dto.UserSignUp) (string, error) {
	hPassword, err := s.Auth.CreateHashedPassword(input.Password)

	if err != nil {
		return "", err
	}

	user, _ := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hPassword,
		Phone:    input.Phone,
	})

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s *UserService) findUserByEmail(email string) (*domain.User, error) {
	user, err := s.Repo.FindUser(email)

	return &user, err
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, err := s.findUserByEmail(email)
	if err != nil {
		return "", errors.New("user does not exist")
	}

	// Compare password and generate token if successful login
	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	// generate token
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s *UserService) GetVerificationCode(e domain.User) (int, error) {

	return 0, nil
}

func (s *UserService) VerifyCode(id uint, code int) error {

	return nil
}

func (s *UserService) CreateProfile(id uint, input any) error {

	return nil
}

func (s *UserService) GetProfile(id uint) (*domain.User, error) {

	return nil, nil
}

func (s *UserService) UpdateProfile(id uint, input any) error {

	return nil
}

func (s *UserService) BecomeSeller(id uint, input any) (string, error) {
	return "", nil
}

func (s *UserService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}

func (s *UserService) CreateCart(input any, u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s *UserService) CreateOrder(u domain.User) (int, error) {
	return 0, nil
}

func (s *UserService) GetOrders(u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s *UserService) GetOrderById(id uint, uId uint) (interface{}, error) {
	return nil, nil
}
