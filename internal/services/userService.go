package services

import (
	"errors"
	"fmt"
	"go-ecommerce-app/configs"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/pkg/notification"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	CRepo  repository.CatalogRepository
	Auth   helper.Auth
	Config configs.AppConfig
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

func (s *UserService) isVerifiedUser(id uint) bool {
	currentUser, err := s.Repo.FindUserByID(id)

	return err == nil && currentUser.Verified
}

func (s *UserService) GetVerificationCode(e domain.User) error {

	// if user already verified
	if s.isVerifiedUser(e.ID) {
		return errors.New("user already verified")
	}

	// generate verification code
	code, err := s.Auth.GenerateCode()
	if err != nil {
		return err
	}

	// update user with verification code
	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}
	_, err = s.Repo.UpdateUser(e.ID, user)
	if err != nil {
		return errors.New("unable to update user with verification code")
	}

	user, _ = s.Repo.FindUserByID(e.ID)
	msg := fmt.Sprintf("Your verification code is %d", code)

	// send SMS
	notificationClient := notification.NewNotificationClient(s.Config)
	err = notificationClient.SendSMS(user.Phone, msg)
	if err != nil {
		return errors.New("unable to send SMS")
	}

	return nil
}

func (s *UserService) VerifyCode(id uint, code int) error {
	// if user already verified
	if s.isVerifiedUser(id) {
		return errors.New("user already verified")
	}

	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return err
	}
	if user.Code != code {
		return errors.New("verification code does not match")
	}
	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code expired")
	}

	updateUser := domain.User{
		Verified: true,
	}
	_, err = s.Repo.UpdateUser(id, updateUser)
	if err != nil {
		return errors.New("unable to verify user")
	}

	return nil
}

func (s *UserService) CreateProfile(id uint, input dto.ProfileInput) error {

	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	_, err = s.Repo.UpdateUser(id, user)
	if err != nil {
		return err
	}

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserID:       id,
	}
	err = s.Repo.CreateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetProfile(id uint) (*domain.User, error) {
	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateProfile(id uint, input dto.ProfileInput) error {
	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = s.Repo.UpdateUser(id, user)
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserID:       id,
	}

	err = s.Repo.UpdateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {
	// Find existing user
	user, _ := s.Repo.FindUserByID(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("user is already a seller")
	}

	// update user
	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})
	if err != nil {
		return "", err
	}

	// generate token
	token, _ := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	// Create bank account
	err = s.Repo.CreateBankAccount(domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		UserId:      id,
		PaymentType: input.PaymentType,
	})

	return token, err
}

func (s *UserService) FindCart(id uint) ([]domain.Cart, error) {
	cartItems, err := s.Repo.FindCartItems(id)

	return cartItems, err
}

func (s *UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {
	// check if the cart exists for the user
	cart, _ := s.Repo.FindCartItem(u.ID, input.ProductId)

	if cart.ID > 0 {
		if input.ProductId == 0 {
			return nil, errors.New("product id is required")
		}

		if input.Qty < 1 {
			err := s.Repo.DeleteCartById(cart.ID)
			if err != nil {
				return nil, errors.New("failed to delete cart item")
			}
		} else {
			cart.Qty = input.Qty
			err := s.Repo.UpdateCart(cart)
			if err != nil {
				return nil, errors.New("failed to update cart item")
			}
		}
	} else {
		product, _ := s.CRepo.FindProductByID(int(input.ProductId))
		if product.ID < 1 {
			return nil, errors.New("product does not exist")
		}

		err := s.Repo.CreateCart(domain.Cart{
			ProductId: input.ProductId,
			UserId:    u.ID,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       input.Qty,
			Price:     product.Price,
			SellerId:  uint(product.UserId),
		})
		if err != nil {
			return nil, errors.New("failed to add product to cart")
		}
	}

	return s.Repo.FindCartItems(u.ID)
}

func (s *UserService) CreateOrder(u domain.User) (int, error) {
	// get cart items for the user
	cartItems, err := s.Repo.FindCartItems(u.ID)
	if err != nil {
		return 0, errors.New("cart does not exist")
	}

	if len(cartItems) == 0 {
		return 0, errors.New("cart is empty. cannot create order")
	}

	// find success payment reference status
	paymentId := "pay_12345"
	txnId := "txn_12345"
	orderRefId, _ := helper.RandomNumbers(8)

	// create order with generated order reference
	var amount float64
	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		amount += item.Price * float64(item.Qty)
		orderItems = append(orderItems, domain.OrderItem{
			ProductId: item.ProductId,
			Qty:       int(item.Qty),
			Price:     item.Price,
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			SellerId:  item.SellerId,
		})
	}

	order := domain.Order{
		UserId:         u.ID,
		PaymentId:      paymentId,
		TransactionId:  txnId,
		OrderRefNumber: uint(orderRefId),
		Items:          orderItems,
	}
	err = s.Repo.CreateOrder(order)

	// send notification to user

	// remove cart items

	// return order reference

	return 0, nil
}

func (s *UserService) GetOrders(u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s *UserService) GetOrderById(id uint, uId uint) (interface{}, error) {
	return nil, nil
}
