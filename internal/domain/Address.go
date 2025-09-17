package domain

import "time"

type Address struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2"`
	City         string    `json:"city"`
	PostCode     string    `json:"postCode"`
	Country      string    `json:"country"`
	UserID       uint      `json:"user_id"`
	CreatedAt    time.Time `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp"`
}
