package models

import "time"

type GormBook struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Author    string `gorm:"not null"`
	Year      int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (GormBook) TableName() string {
	return "gorm_books"
}
