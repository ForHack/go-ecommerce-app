package domain

import "time"

type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id"`
	ProductId uint      `json:"product_id"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
	SellerId  uint      `json:"seller_id"`
	Price     float64   `json:"price"`
	Qty       uint      `json:"qty"`
	CreatedAt time.Time `gorm:"current_timestamp"`
	UpdatedAt time.Time `gorm:"current_timestamp"`
}
