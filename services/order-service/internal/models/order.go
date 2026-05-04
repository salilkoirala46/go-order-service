package models

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint        `gorm:"not null;index"           json:"user_id"`
	Product   string      `gorm:"not null"                 json:"product"`
	Quantity  int         `gorm:"not null"                 json:"quantity"`
	Price     float64     `gorm:"not null"                 json:"price"`
	Total     float64     `gorm:"not null"                 json:"total"`
	Status    OrderStatus `gorm:"default:'pending'"        json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
