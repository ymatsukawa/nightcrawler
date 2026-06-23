package main

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func env(key, val string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return val
}

func dsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		env("DB_HOST", "db"),
		env("DB_USER", "root"),
		env("DB_PASSWORD", "password"),
		env("DB_NAME", "testdb"),
		env("DB_PORT", "5432"),
	)
}

func openAndPing(logger *slogGormLogger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn()), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB(logger *slogGormLogger) (*gorm.DB, error) {
	var lastErr error
	for i := 0; i < 30; i++ {
		db, err := openAndPing(logger)
		if err == nil {
			return db, nil
		}
		lastErr = err
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("could not connect to db after retries: %w", lastErr)
}
