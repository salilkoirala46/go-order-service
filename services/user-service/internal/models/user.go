package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null"                 json:"name"`
	Email     string    `gorm:"uniqueIndex;not null;size:191" json:"email"`
	Password  string    `gorm:"not null"                 json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
