package dto

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignUp struct {
	UserLogin
	Phone string `json:"phone"`
}

type VerificationCodeInput struct {
	Code int `json:"code"`
}

type SellerInput struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	BankAccountNumber uint   `json:"bankAccountNumber"`
	SwiftCode         string `json:"swiftCode"`
	PaymentType       string `json:"paymentType"`
}

type AddressInput struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	PostCode     string `json:"postCode"`
	Country      string `json:"country"`
}

type ProfileInput struct {
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	AddressInput AddressInput `json:"address"`
}
