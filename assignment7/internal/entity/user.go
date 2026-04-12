package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Username string    `json:"username" gorm:"size:100;uniqueIndex;not null"`
	Email    string    `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password string    `json:"-" gorm:"not null"`
	Role     string    `json:"role" gorm:"size:32;not null;default:user"`
	Verified bool      `json:"verified" gorm:"not null;default:false"`
}
