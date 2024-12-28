package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DSN string
}

func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	// Открываем соединение с базой данных
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удалось подключиться
	}

	// Получаем объект *sql.DB из *gorm.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Проверяем соединение с базой данных
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to the PostgreSQL database successfully")
	return db, nil
}
