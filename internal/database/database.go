package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL bulunamadı")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("PostgreSQL bağlantısı başarılı")
	return db, nil
}
