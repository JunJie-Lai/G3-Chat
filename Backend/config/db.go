package config

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	//db.SetMaxOpenConns()
	//db.SetMaxIdleConns()
	//db.SetConnMaxIdleTime()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
