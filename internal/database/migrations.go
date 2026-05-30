package database

import (
	"context"
	"database/sql"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	queries := []string{
		`
	CREATE TABLE IF NOT EXISTS books(
	    id SERIAL PRIMARY KEY,
	    title TEXT NOT NULL,
	    author TEXT NOT NULL,
	    year INT NOT NULL,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
		`,
		`
		CREATE TABLE IF NOT EXISTS users(
		  id SERIAL PRIMARY KEY,
		  name TEXT NOT NULL,
		  email TEXT NOT NULL UNIQUE,
		  password_hash TEXT NOT NULL,
		  role TEXT NOT NULL DEFAULT 'user',
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  
		);
`,
	}

	for _, query := range queries {
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}
