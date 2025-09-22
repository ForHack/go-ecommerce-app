package domain

import "time"

type Payment struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	UserId        uint          `json:"user_id"`
	CaptureMethod string        `json:"capture_method"`
	Amount        float64       `json:"amount"`
	OrderId       string        `json:"order_id"`
	TransactionId uint          `json:"transaction_id"`
	CustomerId    string        `json:"customer_id"`
	PaymentId     string        `json:"payment_id"`
	Status        PaymentStatus `json:"status" gorm:"default:'initial'"`
	Response      string        `json:"response"`
	PaymentUrl    string        `json:"payment_url"`
	CreatedAt     time.Time     `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time     `gorm:"default:current_timestamp"`
}

type PaymentStatus string

const (
	PaymentStatusInitial PaymentStatus = "initial"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
)
