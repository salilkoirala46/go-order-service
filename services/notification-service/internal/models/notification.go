package models

import "time"

type NotificationType string

const (
	TypeWelcome           NotificationType = "WELCOME"
	TypeOrderCreated      NotificationType = "ORDER_CREATED"
	TypeOrderStatusChange NotificationType = "ORDER_STATUS_CHANGE"
)

type Notification struct {
	ID        uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint             `gorm:"not null;index"           json:"user_id"`
	Type      NotificationType `gorm:"not null"                 json:"type"`
	Message   string           `gorm:"not null"                 json:"message"`
	Read      bool             `gorm:"default:false"            json:"read"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
