package gormdb

import (
	"BookHub/internal/gormdb/models"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env dosyası yüklenemedi, ortam değişkenleri kullanılacak")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL bulunamadı")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("GORM PostgreSQL bağlantısı başarılı")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.GormBook{}, &models.GormUser{})
	if err != nil {
		return err
	}

	fmt.Println("GORM migration başarılı")
	return nil
}
