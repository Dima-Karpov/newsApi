package domain

import (
	"github.com/google/uuid"
	"time"
)

// Модели данных

type RSSItem struct {
	Title       string
	Description string
	PublishedAt time.Time
	Link        string
}

type NewsList struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"` // Поле id как SERIAL
	Title       string    `gorm:"size:255;not null"`                    // VARCHAR(255)
	Description string    `gorm:"type:text;not null"`                   // TEXT
	PublishedAt time.Time `gorm:"not null"`                             // TIMESTAMP
	Link        string    `gorm:"size:255;unique;not null"`             // VARCHAR(255) UNIQUE
}

// Указываем имя таблицы для GORM
func (NewsList) TableName() string {
	return "news_list"
}
