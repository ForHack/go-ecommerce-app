package payment

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId uint) (*stripe.CheckoutSession, error)
	GetPaymentStatus(pId string) (*stripe.CheckoutSession, error)
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

// CreatePayment implements PaymentClient.
func (p *payment) CreatePayment(amount float64, userId uint, orderId uint) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := amount * 100

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					UnitAmount: stripe.Int64(int64(amountInCents)),
					Currency:   stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("electronic gadget"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(p.successUrl),
		CancelURL:  stripe.String(p.cancelUrl),
	}

	params.AddMetadata("order_id", fmt.Sprintf("%d", orderId))
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))

	session, err := session.New(params)
	if err != nil {
		return nil, errors.New("failed to create checkout session")
	}
	return session, nil
}

// GetPaymentStatus implements PaymentClient.
func (p *payment) GetPaymentStatus(pId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey
	session, err := session.Get(pId, nil)

	if err != nil {
		return nil, errors.New("failed to retrieve payment status")
	}

	return session, nil
}

func NewPaymentClient(stripeSecretKey, successUrl, cancenUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		cancelUrl:       cancenUrl,
	}
}
