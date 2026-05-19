package models

import "time"

type GormUser struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null;default:user"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (GormUser) TableName() string {
	return "gorm_users"
}
