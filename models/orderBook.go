package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderBook struct {
	gorm.Model
	ID      int
	UserID  int
	HotelID int
	Hotel   Hotel
	Nights  int
	Status  string `gorm:"default:pending"`
	Bill    int

	// Midtrans fields
	MidtransOrderID string
	PaymentURL      string
	PaymentMethod   string
	PaidAt          *time.Time
}
